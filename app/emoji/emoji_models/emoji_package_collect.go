package emoji_models

import "beaver/common/models"

// 用户收藏表情包合集
type EmojiPackageCollect struct {
	models.Model
	UserID    string `json:"userId"`    // 用户ID
	PackageID uint   `json:"packageId"` // 表情包ID
	// 注意：移除外键关联，改用关联查询
}
