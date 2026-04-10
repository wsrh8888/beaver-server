package conn

import (
	"encoding/json"
	"sync"

	ws_response "beaver/app/ws/ws_api/response"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/gorilla/websocket"
)

// Client 单个 WebSocket 连接封装，带独立写互斥锁，解决并发写问题
type Client struct {
	Conn *websocket.Conn
	mu   sync.Mutex
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{Conn: conn}
}

// SafeSend 线程安全发送业务消息（含 content/data 层）
func (c *Client) SafeSend(command wsCommandConst.Command, content type_struct.WsContent) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return ws_response.WsResponse(c.Conn, command, content)
}

// SafeSendControl 线程安全发送控制帧（PING/PONG/ACK，无 content/data 层）
func (c *Client) SafeSendControl(frame type_struct.WsControlFrame) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, err := json.Marshal(frame)
	if err != nil {
		logx.Errorf("序列化控制帧失败: %v", err)
		return err
	}
	return c.Conn.WriteMessage(websocket.TextMessage, data)
}
