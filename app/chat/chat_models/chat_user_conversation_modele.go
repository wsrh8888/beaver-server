package chat_models

import (
	"beaver/common/models"
)

type ChatUserConversationModel struct {
	models.Model
	UserID         string `gorm:"not null"`               // 用户id
	ConversationID string `gorm:"not null"`               // 会话id
	LastMessage    string `gorm:""`                       // 最后一条消息内容
	IsDeleted      bool   `gorm:"not null;default:false"` // 标记用户是否删除会话
	IsPinned       bool   `gorm:"not null;default:false"` // 标记会话是否置顶
}
