package emoji_models

import "beaver/common/models"

// 用户收藏的表情
type EmojiCollectEmoji struct {
	models.Model
	UserID     string `json:"userId"`  // 用户ID
	EmojiID    uint   `json:"emojiId"` // 表情ID
	EmojiModel Emoji  `gorm:"foreignkey:EmojiID;references:ID" json:"-"`
}
