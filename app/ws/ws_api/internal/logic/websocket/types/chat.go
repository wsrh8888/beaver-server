package websocket_types

import (
	"beaver/app/chat/chat_rpc/types/chat_rpc"
)

type SendMsgReq struct {
	UserID         string        `json:"userID"`
	ConversationID string        `json:"conversationId"`
	MessageID      string        `json:"messageId"` // 客户端消息ID
	Msg            *chat_rpc.Msg `json:"msg"`
}
