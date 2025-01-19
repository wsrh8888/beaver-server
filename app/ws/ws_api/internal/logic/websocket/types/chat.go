package websocket_types

import "encoding/json"

type SendMsgReq struct {
	UserID         string          `header:"Beaver-User-Id"`
	ConversationID string          `json:"conversationId"`
	Msg            json.RawMessage `json:"msg"`
}
