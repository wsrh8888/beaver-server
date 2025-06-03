package type_struct

import (
	"beaver/common/wsEnum/wsTypeConst"

	"encoding/json"
)

type WsContent struct {
	Timestamp int64  `json:"timestamp"` //消息发送时间
	MessageID string `json:"messageId"` //客户端消息ID
	Data      WsData `json:"data"`      //消息内容
}

type WsData struct {
	Type           wsTypeConst.Type `json:"type"`           // 消息类型
	ConversationID string           `json:"conversationId"` // 会话ID
	Body           json.RawMessage  `json:"body"`           // 消息内容
}

type WsMessage struct {
	Command string    `json:"command"` //命令
	Content WsContent `json:"content"` //消息内容
}
