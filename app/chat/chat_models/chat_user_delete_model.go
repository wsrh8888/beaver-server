package chat_models

import "beaver/common/models"

// ChatUserDelete 记录用户主动删除的消息（仅自己不可见）
// 对标大厂：不直接修改 ChatMessage 表，而是记录用户的删除行为
type ChatUserDelete struct {
	models.Model
	UserID    string `gorm:"size:64;index:idx_user_msg" json:"userId"`
	MessageID string `gorm:"size:64;index:idx_user_msg" json:"messageId"`
	Version   int64  `gorm:"not null;default:0;index" json:"version"` // 版本号，用于多端同步
}
