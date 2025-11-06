package emoji_models

import (
	"beaver/common/models"
)

// 表情
type Emoji struct {
	models.Model
	FileName string `json:"fileName"`                // 文件Url
	Title    string `json:"title"`                   // 表情名称
	AuthorID string `json:"authorId"`                // 创建者ID（用户ID或官方ID）
	Status   int8   `gorm:"default:1" json:"status"` // 状态：1=正常 2=审核中 3=违规禁用
}

// 1、用户上传突破收藏表情
// 在 Emoji表 -》 EmojiCollectEmoji
