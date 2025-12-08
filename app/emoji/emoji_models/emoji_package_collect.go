package emoji_models

import "beaver/common/models"

// 用户收藏表情包合集
type EmojiPackageCollect struct {
	models.Model
	PackageCollectID string `gorm:"column:package_collect_id;size:64;uniqueIndex" json:"packageCollectId"` // 全局唯一ID
	UserID           string `json:"userId"`                                                                // 用户ID
	PackageID        string `gorm:"size:64;index" json:"packageId"`                                        // 表情包ID
	IsDeleted        bool   `gorm:"default:false;index" json:"isDeleted"`                                  // 是否已删除（软删除）
	Version          int64  `gorm:"not null;default:0;index" json:"version"`                               // 基于userId递增
}
