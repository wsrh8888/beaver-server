package emoji_models

import (
	"beaver/common/models"
)

// 表情
type Emoji struct {
	models.Model
	FileName string `json:"fileName"` // 文件Url
	Title    string `json:"title"`    // 表情名称
	AuthorID string `json:"authorId"` // 创建者ID（用户ID或官方ID）
	// 不再直接关联到单个表情包
}

// 1、用户上传突破收藏表情
// 在 Emoji表 -》 EmojiCollectEmoji
