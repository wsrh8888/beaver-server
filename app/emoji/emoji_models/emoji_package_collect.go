package emoji_models

import "beaver/common/models"

// 用户收藏表情包合集
type EmojiPackageCollect struct {
	models.Model
	UUID      string `gorm:"size:64;unique;index" json:"uuid"`        // 全局唯一标识符，用于前端同步
	UserID    string `json:"userId"`                                  // 用户ID
	PackageID uint   `json:"packageId"`                               // 表情包ID
	Version   int64  `gorm:"not null;default:0;index" json:"version"` // 基于userId递增
}
