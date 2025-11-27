package emoji_models

import "beaver/common/models"

// 用户收藏的表情
type EmojiCollectEmoji struct {
	models.Model
	UUID    string `gorm:"size:64;unique;index" json:"uuid"`        // 全局唯一标识符，用于前端同步
	UserID  string `json:"userId"`                                  // 用户ID
	EmojiID uint   `json:"emojiId"`                                 // 表情ID
	Version int64  `gorm:"not null;default:0;index" json:"version"` //基于userId递增
}
