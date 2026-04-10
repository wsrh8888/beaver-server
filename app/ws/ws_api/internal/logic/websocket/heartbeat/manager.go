package heartbeat

import (
	"context"
	"sync"
	"time"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// HandleClientPing 收到客户端 PING，立即回复 PONG（无状态，echo timestamp）
func HandleClientPing(client *ws_conn.Client, timestamp int64) {
	client.SafeSendControl(type_struct.WsControlFrame{
		Command:   wsCommandConst.PONG,
		Timestamp: timestamp,
	})
}

// Manager 服务端主动心跳管理器
// 定时向客户端发送应用级 PING，并维护协议级 WebSocket ping 帧
// 仅在 chatwebsocketlogic.go 中创建和管理，不参与消息循环
type Manager struct {
	client *ws_conn.Client
	userID string
	svcCtx *svc.ServiceContext
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex // 仅用于协议级 ping（直接写原始帧，不经过 SafeSend）
}

func NewManager(client *ws_conn.Client, userID string, svcCtx *svc.ServiceContext) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		client: client,
		userID: userID,
		svcCtx: svcCtx,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (m *Manager) Start() {
	go m.startProtocolPing()
	go m.startApplicationPing()
}

func (m *Manager) Stop() {
	m.cancel()
}

// startProtocolPing 协议级 ping（原始 WebSocket 帧，gorilla 自动处理 pong 维持连接）
func (m *Manager) startProtocolPing() {
	pingPeriod := time.Duration(m.svcCtx.Config.WebSocket.PingPeriod) * time.Second
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			if err := m.sendProtocolPing(); err != nil {
				logx.Errorf("协议级 ping 失败, 用户: %s, 错误: %v", m.userID, err)
				return
			}
		}
	}
}

// startApplicationPing 应用级 PING，客户端收到后应回复 PONG
func (m *Manager) startApplicationPing() {
	interval := m.getAppPingInterval()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.client.SafeSendControl(type_struct.WsControlFrame{
				Command:   wsCommandConst.PING,
				Timestamp: time.Now().UnixMilli(),
			})
		}
	}
}

func (m *Manager) sendProtocolPing() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	writeWait := time.Duration(m.svcCtx.Config.WebSocket.WriteWait) * time.Second
	if err := m.client.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}
	return m.client.Conn.WriteMessage(websocket.PingMessage, []byte{})
}

func (m *Manager) getAppPingInterval() time.Duration {
	if m.svcCtx.Config.WebSocket.AppHeartbeatInterval > 0 {
		return time.Duration(m.svcCtx.Config.WebSocket.AppHeartbeatInterval) * time.Second
	}
	return time.Duration(m.svcCtx.Config.WebSocket.PingPeriod*2) * time.Second
}
