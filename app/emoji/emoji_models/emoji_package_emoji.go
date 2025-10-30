package emoji_models

import "beaver/common/models"

// EmojiPackageEmoji 表情包与表情的多对多关联表
// 一个表情可以属于多个表情包，一个表情包可以包含多个表情
type EmojiPackageEmoji struct {
	models.Model
	PackageID  uint         `json:"packageId"` // 表情包ID
	EmojiID    uint         `json:"emojiId"`   // 表情ID
	SortOrder  int          `json:"sortOrder"` // 在表情包中的排序
	Package    EmojiPackage `gorm:"foreignkey:PackageID;references:Id" json:"-"`
	EmojiModel Emoji        `gorm:"foreignkey:EmojiID;references:Id" json:"-"`
}
