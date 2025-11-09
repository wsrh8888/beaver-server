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
	ConversationID string `gorm:"size:128;uniqueIndex" json:"conversationId"` // 唯一会话ID（私聊/群聊/系统）
	Type           int    `gorm:"not" json:"type"`                            // 1=私聊 2=群聊 3=系统会话
	MaxSeq         int64  `gorm:"not;default:0" json:"maxSeq"`                // 会话全局最新消息序号
	LastMessage    string `gorm:"size:256" json:"lastMessage"`                // 会话最后一条消息预览（全局唯一）
	Version        int64  `gorm:"not;default:0;index" json:"version"`         // 公共元信息版本，用于同步（基于ConversationID递增，从0开始）
}
