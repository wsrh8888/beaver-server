package moment_models

import (
	"beaver/common/models"
)

/**
 * @description: 评论表
 */
type MomentCommentModel struct {
	models.Model
	UUID      string `gorm:"size:64;uniqueIndex;not null" json:"uuid"`      // 全局唯一ID (UUID，跨库同步用)
	MomentID  string `gorm:"size:64;not null;index" json:"momentId"`        // 动态UUID (索引，关联moment.uuid)
	UserID    string `gorm:"size:64;not null;index" json:"userId"`          // 评论用户Id (索引)
	Content   string `gorm:"type:text;not null" json:"content"`             // 评论内容
	IsDeleted bool   `gorm:"not null;default:false;index" json:"isDeleted"` // 软删除标记 (索引)
}
