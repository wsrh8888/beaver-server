package moment_models

import (
	"beaver/common/models"
)

/**
 * @description: 评论表
 */
type MomentCommentModel struct {
	models.Model
	UUID         string `gorm:"size:64;uniqueIndex;not null" json:"uuid"`      // 全局唯一ID (UUID，跨库同步用)
	MomentID     string `gorm:"size:64;not null;index" json:"momentId"`        // 动态UUID (索引，关联moment.uuid)
	UserID       string `gorm:"size:64;not null;index" json:"userId"`          // 评论用户Id (索引)
	MomentUserID string `gorm:"size:64;not null;index" json:"momentUserId"`    // 动态发布者Id (索引，评论时递增该用户的版本号)
	Content      string `gorm:"type:text;not null" json:"content"`             // 评论内容
	IsDeleted    bool   `gorm:"not null;default:false;index" json:"isDeleted"` // 软删除标记 (索引)
	Version      int64  `gorm:"not null;default:0;index" json:"version"`       // 用户级版本号（基于MomentUserID递增）⭐
}
