package heartbeat

import (
	"context"
	"fmt"
	"time"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/core/coreonline"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// HandleClientPing 收到客户端 PING，续期在线态并回复 PONG。
func HandleClientPing(rdb *redis.Client, userID, slot, instanceID string, client *ws_conn.Client, timestamp int64) {
	coreonline.MarkOnline(rdb, userID, slot, instanceID)
	err := client.SafeSendControl(type_struct.WsControlFrame{
		Command:   wsCommandConst.PONG,
		Timestamp: timestamp,
	})
	if err != nil {
		fmt.Printf("回复 PONG 失败: 错误: %v, 时间戳: %d\n", err, timestamp)
	} else {
		fmt.Printf("已回复 PONG: 时间戳: %d\n", timestamp)
	}
}

// Manager 服务端主动心跳管理器
// 维护协议级 WebSocket ping 帧，确保链路不被中间件（如 Nginx/ELB）超时断开
type Manager struct {
	client      *ws_conn.Client
	userID      string
	deviceGroup string
	svcCtx      *svc.ServiceContext
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewManager(client *ws_conn.Client, userID, deviceGroup string, svcCtx *svc.ServiceContext) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		client:      client,
		userID:      userID,
		deviceGroup: deviceGroup,
		svcCtx:      svcCtx,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (m *Manager) Start() {
	go m.startProtocolPing()
	coreonline.MarkOnline(m.svcCtx.Redis, m.userID, m.deviceGroup, m.svcCtx.InstanceID)
}

func (m *Manager) Stop() {
	m.cancel()
}

// startProtocolPing 协议级 ping（原始 WebSocket 帧）
func (m *Manager) startProtocolPing() {
	pingPeriod := time.Duration(m.svcCtx.Config.WebSocket.PingPeriod) * time.Second
	if pingPeriod <= 0 {
		pingPeriod = 30 * time.Second
	}
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			coreonline.MarkOnline(m.svcCtx.Redis, m.userID, m.deviceGroup, m.svcCtx.InstanceID)
			if err := m.sendProtocolPing(); err != nil {
				logx.Errorf("协议级 ping 失败, 用户: %s, 错误: %v", m.userID, err)
				return
			}
		}
	}
}

func (m *Manager) sendProtocolPing() error {
	// 重要：使用 Client 统一的互斥锁，严禁并发写 Conn
	m.client.Mu.Lock()
	defer m.client.Mu.Unlock()

	writeWait := time.Duration(m.svcCtx.Config.WebSocket.WriteWait) * time.Second
	if writeWait <= 0 {
		writeWait = 10 * time.Second
	}

	if err := m.client.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}
	return m.client.Conn.WriteMessage(websocket.PingMessage, []byte{})
}
