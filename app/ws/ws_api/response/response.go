package ws_response

import (
	type_struct "beaver/app/ws/ws_api/types"
	utils "beaver/utils/rand"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Response struct {
	Code       int                   `json:"code"`
	Command    string                `json:"command"`
	Content    type_struct.WsContent `json:"content"`
	MessageID  string                `json:"messageId"`
	ServerTime int64                 `json:"serverTime"`
}

func WsResponse(conn *websocket.Conn, command string, content type_struct.WsContent) {
	code := 0

	response := Response{
		Command:    command,
		Code:       code,
		Content:    content,
		MessageID:  utils.GenerateRandomString(8),
		ServerTime: time.Now().Unix(),
	}
	responseJSON, _ := json.Marshal(response)
	conn.WriteMessage(websocket.TextMessage, responseJSON)
}
