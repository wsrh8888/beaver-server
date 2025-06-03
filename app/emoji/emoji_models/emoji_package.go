package emoji_models

import "beaver/common/models"

type EmojiPackage struct {
	models.Model
	Title       string `json:"title"`       // 表情包名称
	CoverFile   string `json:"coverFile"`   // 表情包封面文件
	UserID      string `json:"userID"`      // 创建者ID（用户ID或官方ID）
	Description string `json:"description"` // 表情包描述
	Type        string `json:"type"`        // 类型：official-官方，user-用户自定义
	Status      int    `json:"status"`      // 状态：1-正常，0-禁用
	// 不再直接关联表情
}
