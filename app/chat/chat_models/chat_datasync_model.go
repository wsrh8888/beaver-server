package chat_models

import (
	"beaver/common/models"
)

// 存储会话级别的信息（所有用户共享）
// 记录会话类型（私聊/群聊）
// 记录会话的最新消息序列号
// 用于数据同步（客户端需要知道哪些会话有更新）

// ChatConversationMeta 数据同步模型
type ChatConversationMeta struct {
	models.Model
	ConversationID string `json:"conversationId"`      // 会话id（单聊为用户id，群聊为群id）
	Type           int    `gorm:"not"`                 // 1=私聊 2=群聊 3=系统会话/客服
	LastReadSeq    int64  `gorm:"not;default:0"`       // 会话消息的最大 Seq（用于消息定位）
	Version        int64  `gorm:"not;default:0;index"` // 会话元信息版本（类型、参与人等有变时+1）
}
