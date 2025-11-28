package emoji_models

import "beaver/common/models"

// EmojiPackageEmoji 表情包与表情的多对多关联表
// 一个表情可以属于多个表情包，一个表情包可以包含多个表情
type EmojiPackageEmoji struct {
	models.Model
	UUID      string `gorm:"size:64;unique;index" json:"uuid"`        // 全局唯一标识符，用于前端同步
	PackageID string `gorm:"size:64;index" json:"packageId"`          // 表情包UUID
	EmojiID   string `gorm:"size:64;index" json:"emojiId"`            // 表情UUID
	SortOrder int    `gorm:"default:0" json:"sortOrder"`              // 在表情包中的排序
	Version   int64  `gorm:"not null;default:0;index" json:"version"` // 基于PackageID递增的内容版本号
}
