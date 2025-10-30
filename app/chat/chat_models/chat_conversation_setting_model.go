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
	UserID         string `gorm:"size:64;not;index"`   // 属于哪个用户
	ConversationID string `json:"conversationId"`      // 会话id（单聊为用户id，群聊为群id）
	LastMessage    string `gorm:"size:128"`            // 最后一条消息预览（UI显示用）
	IsDeleted      bool   `gorm:"default:false"`       // 用户是否把会话从列表删除/隐藏
	IsPinned       bool   `gorm:"default:false"`       // 用户置顶
	IsMuted        bool   `gorm:"default:false"`       // 用户免打扰
	LastReadSeq    int64  `gorm:"not;default:0"`       // 用户已读到的消息 Seq（用户维度的已读）
	Version        int64  `gorm:"not;default:0;index"` // 用户个性化设置版本
}
