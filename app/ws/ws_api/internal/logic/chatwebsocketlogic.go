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
	// 1. 升级 HTTP → WebSocket
	conn, err := upgradeToWebSocket(w, r)
	if err != nil {
		logx.Errorf("WebSocket升级失败, 用户: %s, 错误: %v", req.UserID, err)
		return nil, nil
	}

	// 2. 鉴权：JWT 签名 + Redis 登录态
	if authErr := ws_auth.VerifyWsToken(req.Token, l.svcCtx.Config.Auth.AccessSecret, req.UserID, l.svcCtx.Redis); authErr != nil {
		logx.Errorf("WS鉴权失败, 用户: %s, 错误: %v", req.UserID, authErr)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, authErr.Error()))
		conn.Close()
		return nil, nil
	}

	// 3. 配置连接参数
	configureWebSocketConn(conn, l.svcCtx)

	// 4. 封装为 Client（带写 mutex）
	client := ws_conn.NewClient(conn)

	// 5. 注册连接，从 User-Agent 识别设备类型
	userAgent := r.Header.Get("User-Agent")
	deviceType := device.GetDeviceType(userAgent)
	userKey := ws_conn.GetUserKey(req.UserID, deviceType)

	logx.Infof("用户上线: %s, 设备: %s, 地址: %s", req.UserID, deviceType, conn.RemoteAddr().String())
	manageUserConnection(userKey, client, req.UserID, deviceType)
	defer cleanupConnection(req.UserID, conn)

	// 6. 启动心跳
	heartbeatManager := heartbeat.NewManager(client, req.UserID, l.svcCtx)
	defer heartbeatManager.Stop()
	heartbeatManager.Start()

	// 7. 消息循环（阻塞直到连接断开）
	ws.HandleWebSocketMessages(l.ctx, l.svcCtx, req, r, client)

	return nil, nil
}

func manageUserConnection(userKey string, client *ws_conn.Client, userID, deviceType string) {
	ws_conn.WsMapMutex.Lock()
	defer ws_conn.WsMapMutex.Unlock()

	addr := client.Conn.RemoteAddr().String()
	userWsInfo, ok := ws_conn.UserOnlineWsMap[userKey]

	if ok {
		// desktop/mobile 限制单连接，关闭旧连接
		if deviceType == "desktop" || deviceType == "mobile" {
			for oldAddr, oldClient := range userWsInfo.WsClientMap {
				logx.Infof("关闭旧连接, 用户: %s, 设备: %s, 地址: %s", userID, deviceType, oldAddr)
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

	logx.Infof("连接注册成功, 用户: %s, 设备: %s", userID, deviceType)
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

func cleanupConnection(userID string, conn *websocket.Conn) {
	conn.Close()
	addr := conn.RemoteAddr().String()

	logx.Infof("清理连接, 用户: %s, 地址: %s", userID, addr)

	ws_conn.WsMapMutex.Lock()
	defer ws_conn.WsMapMutex.Unlock()

	for _, dt := range []string{"desktop", "mobile", "web", "unknown"} {
		userKey := ws_conn.GetUserKey(userID, dt)
		userWsInfo, ok := ws_conn.UserOnlineWsMap[userKey]
		if !ok {
			continue
		}
		delete(userWsInfo.WsClientMap, addr)
		if len(userWsInfo.WsClientMap) == 0 {
			delete(ws_conn.UserOnlineWsMap, userKey)
			logx.Infof("用户下线, 用户: %s, 设备: %s", userID, dt)
		}
	}
}
