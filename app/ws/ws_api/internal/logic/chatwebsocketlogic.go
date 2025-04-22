// Package logic 实现WebSocket聊天服务的业务逻辑
package logic

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	ws "beaver/app/ws/ws_api/internal/logic/websocket"
	websocket_utils "beaver/app/ws/ws_api/internal/logic/websocket/utils"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"

	"github.com/gorilla/websocket"

	"github.com/zeromicro/go-zero/core/logx"
)

// wsMapMutex 用于保护UserOnlineWsMap的互斥锁
// 确保多个goroutine同时访问用户连接映射时的并发安全
var wsMapMutex sync.RWMutex

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
		return
	}

	// 配置WebSocket连接参数（超时、最大消息大小等）
	configureWebSocketConn(conn, l.svcCtx)

	// 当函数返回时关闭连接并清理资源
	defer cleanupConnection(req.UserID, conn)

	fmt.Println("用户上线", req.UserID)

	// 使用读锁查询用户连接信息
	wsMapMutex.RLock()
	userWsInfo, ok := websocket_utils.UserOnlineWsMap[req.UserID]
	wsMapMutex.RUnlock()

	// 获取连接的远程地址，用作连接标识
	addr := conn.RemoteAddr().String()

	if !ok {
		// 用户首次连接，创建新的用户连接信息
		userWsInfo = &websocket_utils.UserWsInfo{
			WsClientMap: map[string]*websocket.Conn{
				addr: conn,
			},
		}

		// 使用写锁添加用户连接信息
		wsMapMutex.Lock()
		websocket_utils.UserOnlineWsMap[req.UserID] = userWsInfo
		wsMapMutex.Unlock()
	} else {
		// 用户已有其他连接，添加新连接到连接映射
		wsMapMutex.Lock()
		userWsInfo.WsClientMap[addr] = conn
		wsMapMutex.Unlock()
	}

	// 启动单独的goroutine处理心跳检测
	go startPingPong(conn, l.svcCtx)

	// 处理WebSocket消息
	ws.HandleWebSocketMessages(l.ctx, l.svcCtx, req, r, conn)
	return &types.WsRes{}, nil
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

// startPingPong 启动心跳检测
// 定期向客户端发送ping消息，确保连接活跃
// conn: 需要进行心跳检测的WebSocket连接
// svcCtx: 服务上下文，包含配置信息
func startPingPong(conn *websocket.Conn, svcCtx *svc.ServiceContext) {
	// 从配置中获取心跳间隔时间
	pingPeriod := time.Duration(svcCtx.Config.WebSocket.PingPeriod) * time.Second

	// 创建定时器，按pingPeriod间隔触发
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		// 函数结束时停止定时器并关闭连接
		ticker.Stop()
		conn.Close()
	}()

	// 写入超时时间
	writeWait := time.Duration(svcCtx.Config.WebSocket.WriteWait) * time.Second

	// 循环发送心跳
	for {
		select {
		case <-ticker.C:
			// 设置写入超时
			if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}

			// // 发送ping消息
			// if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			// 	logx.Error("发送心跳失败:", err)
			// 	return // 发送失败则退出心跳goroutine
			// }
		}
	}
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
// 在连接关闭时移除用户连接映射
// userID: 用户ID
// conn: 需要清理的WebSocket连接
func cleanupConnection(userID string, conn *websocket.Conn) {
	// 关闭连接
	conn.Close()
	// 获取连接地址
	addr := conn.RemoteAddr().String()

	// 使用互斥锁保护共享资源操作
	wsMapMutex.Lock()
	defer wsMapMutex.Unlock()

	// 从用户连接映射中删除此连接
	userWsInfo, ok := websocket_utils.UserOnlineWsMap[userID]
	if ok {
		delete(userWsInfo.WsClientMap, addr)
	}

	// 如果用户没有任何活跃连接，从用户映射中删除用户
	if userWsInfo != nil && len(userWsInfo.WsClientMap) == 0 {
		delete(websocket_utils.UserOnlineWsMap, userID)
	}
}
