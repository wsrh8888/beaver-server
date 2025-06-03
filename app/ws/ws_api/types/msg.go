package type_struct

import "encoding/json"

type BodySendMsg struct {
	ConversationID string          `json:"conversationId"`
	MessageID      string          `json:"messageId"` // 客户端消息ID
	Msg            json.RawMessage `json:"msg"`
}
