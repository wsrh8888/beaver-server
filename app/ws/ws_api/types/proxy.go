package type_struct

import "encoding/json"

type BodyProxyMsg struct {
	ConversationID string          `json:"conversationId"`
	Content        json.RawMessage `json:"content"`
}
