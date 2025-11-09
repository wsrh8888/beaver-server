package chat_models

import (
	"beaver/common/models"
)

// 存储用户级别的会话设置（每个用户独立）
// 记录用户对会话的个性化操作（置顶、免打扰、删除等）
// 记录用户在该会话中的已读状态
// 用于UI显示（会话列表、未读消息等）

type ChatUserConversation struct {
	models.Model
	UserID         string `gorm:"size:64;index" json:"userId"`          // 用户ID
	ConversationID string `gorm:"size:128;index" json:"conversationId"` // 关联的会话ID
	IsHidden       bool   `gorm:"default:false" json:"isHidden"`        // 是否在当前用户的会话列表隐藏
	IsPinned       bool   `gorm:"default:false" json:"isPinned"`        // 置顶
	IsMuted        bool   `gorm:"default:false" json:"isMuted"`         // 免打扰
	UserReadSeq    int64  `gorm:"not;default:0" json:"userReadSeq"`     // 当前用户已读游标
	Version        int64  `gorm:"not;default:0;index" json:"version"`   // 配置版本，用于多端同步（基于ConversationID递增，从0开始）
}
