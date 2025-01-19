package emoji_models

import "beaver/common/models"

// 用户收藏表情包合集
type EmojiPackageCollect struct {
	models.Model
	UserID    string       `json:"userId"`    // 用户ID
	PackageID uint         `json:"packageId"` // 表情包ID
	Package   EmojiPackage `gorm:"foreignkey:PackageID;references:ID" json:"-"`
}
