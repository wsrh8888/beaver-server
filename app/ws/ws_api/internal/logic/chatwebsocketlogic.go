package logic

import (
	"context"
	"net/http"
	"time"

	ws "beaver/app/ws/ws_api/internal/logic/websocket"
	ws_auth "beaver/app/ws/ws_api/internal/logic/websocket/auth"
	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/logic/websocket/heartbeat"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	"beaver/core/corepush"
	"beaver/utils/device"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChatWebsocketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatWebsocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatWebsocketLogic {
	return &ChatWebsocketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *ChatWebsocketLogic) ChatWebsocket(req *types.WsReq, w http.ResponseWriter, r *http.Request) (resp *types.WsRes, err error) {
	// 1. 基础入口校验：精准识别设备
	userAgent := r.Header.Get("User-Agent")
	preciseType := device.GetDeviceType(userAgent)
	if preciseType == device.DeviceUnknown {
		logx.Errorf("连接拒绝：非法设备接入, 用户: %s, UA: %s", req.UserID, userAgent)
		http.Error(w, "Illegal Device", http.StatusForbidden)
		return nil, nil
	}

	// 2. 获取所属槽位 (Group: mobile/desktop)，用于检索 Redis 登录态
	deviceGroup := device.GetDeviceGroup(preciseType)

	// 3. 鉴权：先于升级执行。使用槽位 (Group) 进行登录态比对
	if authErr := ws_auth.VerifyWsToken(req.Token, l.svcCtx.Config.Auth.AccessSecret, req.UserID, deviceGroup, l.svcCtx.Redis); authErr != nil {
		logx.Errorf("WS鉴权失败, 用户: %s, 精准设备: %s, 槽位: %s, 错误: %v", req.UserID, preciseType, deviceGroup, authErr)
		http.Error(w, authErr.Error(), http.StatusUnauthorized)
		return nil, nil
	}

	// 2. 升级 HTTP → WebSocket
	conn, err := upgradeToWebSocket(w, r)
	if err != nil {
		logx.Errorf("WebSocket升级失败, 用户: %s, 错误: %v", req.UserID, err)
		return nil, nil
	}

	// 3. 配置连接参数
	configureWebSocketConn(conn, l.svcCtx)

	// 4. 封装为 Client（带写 mutex）
	client := ws_conn.NewClient(conn)

	// 5. 注册连接：统一使用槽位 (Group) 管理，确保互踢逻辑闭合
	// 无论具体的 OS 是 windows、macos 还是 linux，在 WS 路由层统一视为 desktop 槽位
	userKey := ws_conn.GetUserKey(req.UserID, deviceGroup)

	logx.Infof("用户上线: %s, 槽位: %s (%s), 地址: %s", req.UserID, deviceGroup, preciseType, conn.RemoteAddr().String())
	manageUserConnection(userKey, client, req.UserID, deviceGroup)
	corepush.MarkOnline(l.svcCtx.Redis, req.UserID, deviceGroup, l.svcCtx.InstanceID)
	connAddr := conn.RemoteAddr().String()
	defer func() {
		conn.Close()
		cleanupConnection(req.UserID, deviceGroup, connAddr, l.svcCtx)
	}()

	// 6. 启动心跳
	heartbeatManager := heartbeat.NewManager(client, req.UserID, deviceGroup, l.svcCtx)
	defer heartbeatManager.Stop()
	heartbeatManager.Start()

	// 7. 消息循环
	ws.HandleWebSocketMessages(l.ctx, l.svcCtx, req, r, client)

	return nil, nil
}

func manageUserConnection(userKey string, client *ws_conn.Client, userID, deviceGroup string) {
	ws_conn.WsMapMutex.Lock()
	defer ws_conn.WsMapMutex.Unlock()

	addr := client.Conn.RemoteAddr().String()
	userWsInfo, ok := ws_conn.UserOnlineWsMap[userKey]

	if ok {
		// 槽位级互踢：目前 desktop 和 mobile 均限制单物理设备在线
		if deviceGroup == "desktop" || deviceGroup == "mobile" {
			for oldAddr, oldClient := range userWsInfo.WsClientMap {
				logx.Infof("【槽位互踢】关闭旧连接, 用户: %s, 槽位: %s, 地址: %s", userID, deviceGroup, oldAddr)
				oldClient.Conn.Close()
				delete(userWsInfo.WsClientMap, oldAddr)
			}
		}
		userWsInfo.WsClientMap[addr] = client
	} else {
		ws_conn.UserOnlineWsMap[userKey] = &ws_conn.UserWsInfo{
			WsClientMap: map[string]*ws_conn.Client{addr: client},
		}
	}

	logx.Infof("连接注册成功, 用户: %s, 槽位: %s", userID, deviceGroup)
}

func configureWebSocketConn(conn *websocket.Conn, svcCtx *svc.ServiceContext) {
	conn.SetReadLimit(int64(svcCtx.Config.WebSocket.MaxMessageSize))
	pongWait := time.Duration(svcCtx.Config.WebSocket.PongWait) * time.Second
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

func upgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	return upGrader.Upgrade(w, r, nil)
}

func cleanupConnection(userID, deviceGroup, addr string, svcCtx *svc.ServiceContext) {
	ws_conn.WsMapMutex.Lock()
	defer ws_conn.WsMapMutex.Unlock()

	userKey := ws_conn.GetUserKey(userID, deviceGroup)
	userWsInfo, ok := ws_conn.UserOnlineWsMap[userKey]
	if !ok {
		corepush.MarkOffline(svcCtx.Redis, userID, deviceGroup, svcCtx.InstanceID)
		return
	}

	delete(userWsInfo.WsClientMap, addr)
	if len(userWsInfo.WsClientMap) == 0 {
		delete(ws_conn.UserOnlineWsMap, userKey)
		corepush.MarkOffline(svcCtx.Redis, userID, deviceGroup, svcCtx.InstanceID)
		logx.Infof("用户下线, 用户: %s, 槽位: %s", userID, deviceGroup)
	}
}
