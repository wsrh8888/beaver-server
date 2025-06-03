package heartbeat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"beaver/app/ws/ws_api/internal/svc"
	ws_response "beaver/app/ws/ws_api/response"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// Manager 心跳管理器
type Manager struct {
	conn   *websocket.Conn
	userID string
	svcCtx *svc.ServiceContext
	ctx    context.Context
	cancel context.CancelFunc
}

// NewManager 创建心跳管理器
func NewManager(conn *websocket.Conn, userID string, svcCtx *svc.ServiceContext) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		conn:   conn,
		userID: userID,
		svcCtx: svcCtx,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start 启动心跳管理
func (m *Manager) Start() {
	go m.startProtocolHeartbeat()
	go m.startApplicationHeartbeat()
}

// Stop 停止心跳管理
func (m *Manager) Stop() {
	m.cancel()
}

// HandleClientHeartbeat 处理客户端心跳
func (m *Manager) HandleClientHeartbeat(content type_struct.WsContent) {
	logx.Infof("收到客户端心跳, 用户: %s", m.userID)

	responseContent := type_struct.WsContent{
		Timestamp: time.Now().UnixMilli(),
		MessageID: m.generateMessageID("heartbeat_response"),
		Data: type_struct.WsData{
			Type: "heartbeat_response",
			Body: json.RawMessage(fmt.Sprintf(`{"server_time": %d}`, time.Now().UnixMilli())),
		},
	}

	ws_response.WsResponse(m.conn, wsCommandConst.HEARTBEAT, responseContent)
	logx.Infof("💗 心跳响应发送成功, 用户: %s", m.userID)
}

// startProtocolHeartbeat 启动协议级心跳
func (m *Manager) startProtocolHeartbeat() {
	pingPeriod := time.Duration(m.svcCtx.Config.WebSocket.PingPeriod) * time.Second
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	logx.Infof("启动协议级心跳, 用户: %s, 间隔: %v", m.userID, pingPeriod)

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			if err := m.sendProtocolPing(); err != nil {
				logx.Errorf("协议级心跳失败, 用户: %s, 错误: %v", m.userID, err)
				return
			}
		}
	}
}

// startApplicationHeartbeat 启动应用级心跳
func (m *Manager) startApplicationHeartbeat() {
	interval := m.getAppHeartbeatInterval()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logx.Infof("启动应用级心跳, 用户: %s, 间隔: %v", m.userID, interval)

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.sendApplicationHeartbeat()
		}
	}
}

// sendProtocolPing 发送协议级ping
func (m *Manager) sendProtocolPing() error {
	writeWait := time.Duration(m.svcCtx.Config.WebSocket.WriteWait) * time.Second
	if err := m.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}

	if err := m.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
		return err
	}

	logx.Debugf("💗 协议级心跳成功, 用户: %s", m.userID)
	return nil
}

// sendApplicationHeartbeat 发送应用级心跳
func (m *Manager) sendApplicationHeartbeat() {
	content := type_struct.WsContent{
		Timestamp: time.Now().UnixMilli(),
		MessageID: m.generateMessageID("server_heartbeat"),
		Data: type_struct.WsData{
			Type: "server_heartbeat",
			Body: json.RawMessage(fmt.Sprintf(`{"server_time": %d}`, time.Now().UnixMilli())),
		},
	}

	ws_response.WsResponse(m.conn, wsCommandConst.HEARTBEAT, content)
	logx.Infof("💓 应用级心跳成功, 用户: %s", m.userID)
}

// getAppHeartbeatInterval 获取应用级心跳间隔
func (m *Manager) getAppHeartbeatInterval() time.Duration {
	if m.svcCtx.Config.WebSocket.AppHeartbeatInterval > 0 {
		return time.Duration(m.svcCtx.Config.WebSocket.AppHeartbeatInterval) * time.Second
	}
	return time.Duration(m.svcCtx.Config.WebSocket.PingPeriod*2) * time.Second
}

// generateMessageID 生成消息ID
func (m *Manager) generateMessageID(prefix string) string {
	return fmt.Sprintf("%s_%d_%s", prefix, time.Now().UnixNano(), m.userID)
}
