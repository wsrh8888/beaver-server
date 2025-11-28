package emoji_models

import "beaver/common/models"

// 用户收藏的表情
type EmojiCollectEmoji struct {
	models.Model
	UUID      string `gorm:"size:64;unique;index" json:"uuid"`        // 全局唯一标识符，用于前端同步
	UserID    string `json:"userId"`                                  // 用户ID
	EmojiID   string `gorm:"size:64;index" json:"emojiId"`            // 表情UUID
	IsDeleted bool   `gorm:"default:false;index" json:"isDeleted"`    // 是否已删除（软删除）
	Version   int64  `gorm:"not null;default:0;index" json:"version"` //基于userId递增
	
}
