package ws_response

import (
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"
	utils "beaver/utils/rand"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Response struct {
	Code       int                    `json:"code"`
	Command    wsCommandConst.Command `json:"command"`
	Content    type_struct.WsContent  `json:"content"`
	MessageID  string                 `json:"messageId"`
	ServerTime int64                  `json:"serverTime"`
}

func WsResponse(conn *websocket.Conn, command wsCommandConst.Command, content type_struct.WsContent) error {
	code := 0

	response := Response{
		Command:    command,
		Code:       code,
		Content:    content,
		MessageID:  utils.GenerateRandomString(8),
		ServerTime: time.Now().Unix(),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		logx.Errorf("序列化WebSocket响应失败: %v", err)
		return err
	}

	// 设置写入超时
	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		logx.Errorf("设置WebSocket写入超时失败: %v", err)
		return err
	}

	if err := conn.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		logx.Errorf("发送WebSocket消息失败: %v", err)
		return err
	}

	return nil
}
