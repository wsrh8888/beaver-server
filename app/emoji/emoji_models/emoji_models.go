package emoji_models

import (
	"beaver/common/models"
)

// 表情
type Emoji struct {
	models.Model
	FileUrl   string       `json:"fileUrl"`   // 文件Url
	Title     string       `json:"title"`     // 表情名称
	PackageID *uint        `json:"packageId"` // 所属表情包ID
	AuthorID  string       `json:"authorId"`  // 创建者ID（用户ID或官方ID）
	Package   EmojiPackage `gorm:"foreignkey:PackageID;references:ID" json:"-"`
}
