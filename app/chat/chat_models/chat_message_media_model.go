package chat_models

import "beaver/common/models"

// ChatMessageMedia 用户对聊天消息的媒体状态（语音已听等）
type ChatMessageMedia struct {
	models.Model
	UserID    string `gorm:"size:64;uniqueIndex:idx_user_msg" json:"userId"`
	MessageID string `gorm:"size:64;uniqueIndex:idx_user_msg" json:"messageId"`
	Version   int64  `gorm:"not null;default:0;index" json:"version"`
}
