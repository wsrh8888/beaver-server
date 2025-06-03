package websocket_types

import "encoding/json"

type SendMsgReq struct {
	UserID         string          `header:"Beaver-User-Id"`
	ConversationID string          `json:"conversationId"`
	MessageID      string          `json:"messageId"` // 客户端消息ID
	Msg            json.RawMessage `json:"msg"`
}
