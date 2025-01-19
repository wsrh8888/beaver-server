package type_struct

import "encoding/json"

type BodySendMsg struct {
	ConversationID string          `json:"conversationId"`
	Msg            json.RawMessage `json:"msg"`
}
