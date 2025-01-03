package type_struct

import "encoding/json"

type BodySendMsg struct {
	ConversationId string          `json:"conversationId"`
	Msg            json.RawMessage `json:"msg"`
}
