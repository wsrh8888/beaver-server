package moment_models

import (
	"beaver/common/models"
)

/**
 * @description: 点赞表
 */
type MomentLikeModel struct {
	models.Model
	UUID         string `gorm:"size:64;uniqueIndex;not null" json:"uuid"`      // 全局唯一ID (UUID，跨库同步用)
	MomentID     uint   `gorm:"not null;index" json:"momentId"`                // 动态Id (索引)
	UserID       string `gorm:"size:64;not null;index" json:"userId"`          // 点赞用户Id (索引)
	MomentUserID string `gorm:"size:64;not null;index" json:"momentUserId"`    // 动态发布者Id (索引，用于版本号递增)
	IsDeleted    bool   `gorm:"not null;default:false;index" json:"isDeleted"` // 软删除标记 (索引)
	Version      int64  `gorm:"not null;default:0;index" json:"version"`       // 动态发布者级版本号（基于MomentUserID递增，每次数据变更都递增）
}
