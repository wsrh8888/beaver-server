// Package logic 实现WebSocket聊天服务的业务逻辑
package logic

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	ws "beaver/app/ws/ws_api/internal/logic/websocket"
	"beaver/app/ws/ws_api/internal/logic/websocket/heartbeat"
	websocket_utils "beaver/app/ws/ws_api/internal/logic/websocket/utils"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// ChatWebsocketLogic WebSocket聊天逻辑的结构体
// 实现WebSocket连接的建立、消息处理与管理
type ChatWebsocketLogic struct {
	logx.Logger                     // 内嵌日志组件
	ctx         context.Context     // 上下文，用于控制请求的生命周期
	svcCtx      *svc.ServiceContext // 服务上下文，包含配置和依赖
}

// NewChatWebsocketLogic 创建一个新的ChatWebsocketLogic实例
// ctx: 请求上下文
// svcCtx: 服务上下文
// 返回: *ChatWebsocketLogic 聊天WebSocket逻辑实例
func NewChatWebsocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatWebsocketLogic {
	return &ChatWebsocketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChatWebsocket 处理WebSocket连接请求的主函数
// req: WebSocket请求参数，包含用户ID等信息
// w: HTTP响应写入器
// r: HTTP请求
// 返回:
//   - resp: WebSocket响应
//   - err: 错误信息
func (l *ChatWebsocketLogic) ChatWebsocket(req *types.WsReq, w http.ResponseWriter, r *http.Request) (resp *types.WsRes, err error) {
	// 将HTTP连接升级到WebSocket连接
	conn, err := UpgradeToWebSocket(w, r)
	if err != nil {
		logx.Errorf("WebSocket升级失败, 用户: %s, 错误: %v", req.UserID, err)
		return
	}

	// 配置WebSocket连接参数
	configureWebSocketConn(conn, l.svcCtx)

	// 创建心跳管理器
	heartbeatManager := heartbeat.NewManager(conn, req.UserID, l.svcCtx)
	defer heartbeatManager.Stop()

	// 当函数返回时关闭连接并清理资源
	defer cleanupConnection(req.UserID, conn)

	// 从User-Agent获取设备类型
	userAgent := r.Header.Get("User-Agent")
	logx.Infof("User-Agent: %s", userAgent)
	deviceType := getDeviceType(userAgent)
	userKey := websocket_utils.GetUserKey(req.UserID, deviceType)

	logx.Infof("用户上线: %s, 设备类型: %s, User-Agent: %s, 远程地址: %s", req.UserID, deviceType, userAgent, conn.RemoteAddr().String())

	// 管理连接映射
	manageUserConnection(userKey, conn, req.UserID, deviceType)

	// 启动心跳管理
	heartbeatManager.Start()

	// 处理WebSocket消息
	ws.HandleWebSocketMessages(l.ctx, l.svcCtx, req, r, conn, heartbeatManager)
	return &types.WsRes{}, nil
}

// manageUserConnection 管理用户连接
func manageUserConnection(userKey string, conn *websocket.Conn, userID, deviceType string) {
	websocket_utils.WsMapMutex.Lock()
	defer websocket_utils.WsMapMutex.Unlock()

	addr := conn.RemoteAddr().String()
	userWsInfo, ok := websocket_utils.UserOnlineWsMap[userKey]

	if !ok {
		userWsInfo = &websocket_utils.UserWsInfo{
			WsClientMap: map[string]*websocket.Conn{addr: conn},
		}
		websocket_utils.UserOnlineWsMap[userKey] = userWsInfo
		logx.Infof("创建用户连接映射, 用户: %s, 设备: %s, userKey: %s", userID, deviceType, userKey)
	} else {
		userWsInfo.WsClientMap[addr] = conn
		logx.Infof("添加连接映射, 用户: %s, 设备: %s, userKey: %s, 连接数: %d", userID, deviceType, userKey, len(userWsInfo.WsClientMap))
	}

	// 打印当前用户的所有连接状态
	logx.Infof("=== 用户连接状态 用户ID: %s ===", userID)
	deviceTypes := []string{"mobile", "windows", "mac", "linux", "web"}
	for _, dt := range deviceTypes {
		uk := websocket_utils.GetUserKey(userID, dt)
		if info, exists := websocket_utils.UserOnlineWsMap[uk]; exists {
			logx.Infof("设备类型: %s, userKey: %s, 连接数: %d", dt, uk, len(info.WsClientMap))
		}
	}
	logx.Infof("=== 连接状态结束 ===")
}

// configureWebSocketConn 配置WebSocket连接的参数
// conn: 需要配置的WebSocket连接
// svcCtx: 服务上下文，包含配置信息
func configureWebSocketConn(conn *websocket.Conn, svcCtx *svc.ServiceContext) {
	// 设置单个消息的最大大小限制
	conn.SetReadLimit(int64(svcCtx.Config.WebSocket.MaxMessageSize))

	// 设置读取超时时间
	pongWait := time.Duration(svcCtx.Config.WebSocket.PongWait) * time.Second
	conn.SetReadDeadline(time.Now().Add(pongWait))

	// 设置Pong处理函数，当收到客户端的pong响应时重置读取超时
	conn.SetPongHandler(func(string) error {
		// 重置读取截止时间，延长连接生命周期
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

// UpgradeToWebSocket 将HTTP连接升级为WebSocket连接
// w: HTTP响应写入器
// r: HTTP请求
// 返回:
//   - *websocket.Conn: 升级后的WebSocket连接
//   - error: 升级过程中的错误
func UpgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	// 创建WebSocket升级器
	upGrader := websocket.Upgrader{
		// 允许所有来源的跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 执行升级
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// cleanupConnection 清理用户连接资源
func cleanupConnection(userID string, conn *websocket.Conn) {
	// 关闭连接
	conn.Close()
	// 获取连接地址
	addr := conn.RemoteAddr().String()

	logx.Infof("开始清理连接资源, 用户: %s, 地址: %s", userID, addr)

	// 使用互斥锁保护共享资源操作
	websocket_utils.WsMapMutex.Lock()
	defer websocket_utils.WsMapMutex.Unlock()

	// 遍历所有设备类型查找并清理连接
	deviceTypes := []string{"mobile", "windows", "mac", "linux", "web"}
	for _, deviceType := range deviceTypes {
		userKey := websocket_utils.GetUserKey(userID, deviceType)
		userWsInfo, ok := websocket_utils.UserOnlineWsMap[userKey]
		if ok {
			delete(userWsInfo.WsClientMap, addr)
			// 如果用户没有任何活跃连接，从用户映射中删除用户
			if len(userWsInfo.WsClientMap) == 0 {
				delete(websocket_utils.UserOnlineWsMap, userKey)
				logx.Infof("删除用户连接映射, 用户: %s, 设备: %s", userID, deviceType)
			} else {
				logx.Infof("保留用户连接映射, 用户: %s, 设备: %s, 剩余连接数: %d", userID, deviceType, len(userWsInfo.WsClientMap))
			}
		}
	}
}

// getDeviceType 根据User-Agent识别设备类型
func getDeviceType(userAgent string) string {
	fmt.Println("userAgent", userAgent)
	userAgent = strings.ToLower(userAgent)
	return "mobile"
	// 移动设备识别
	if strings.Contains(userAgent, "android") {
		return "mobile"
	} else if strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad") {
		return "mobile"
	} else if strings.Contains(userAgent, "mobile") {
		return "mobile"
	} else if strings.Contains(userAgent, "uniapp") {
		return "mobile"
	} else if strings.Contains(userAgent, "uni-app") {
		return "mobile"
	} else if strings.Contains(userAgent, "uni") {
		return "mobile"
	} else if strings.Contains(userAgent, "app") && (strings.Contains(userAgent, "android") || strings.Contains(userAgent, "ios")) {
		return "mobile"
	}

	// 桌面设备识别
	if strings.Contains(userAgent, "windows") {
		return "windows"
	} else if strings.Contains(userAgent, "macintosh") || strings.Contains(userAgent, "mac os") {
		return "mac"
	} else if strings.Contains(userAgent, "linux") {
		return "linux"
	} else {
		return "web"
	}
}
