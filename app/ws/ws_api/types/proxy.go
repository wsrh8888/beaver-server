package type_struct

import "encoding/json"

type BodyProxyMsg struct {
	ConversationId string          `json:"conversationId"`
	Content        json.RawMessage `json:"content"`
}
