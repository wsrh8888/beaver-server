package emoji_models

import "beaver/common/models"

type EmojiPackage struct {
	models.Model
	Title       string  `json:"title"`       // 表情包名称
	CoverFile   string  `json:"coverFile"`   // 表情包封面文件
	UserID      string  `json:"userID"`      // 创建者ID（用户ID或官方ID）
	Description string  `json:"description"` // 表情包描述
	EmojiModels []Emoji `gorm:"foreignkey:PackageID;references:ID" json:"-"`
}
