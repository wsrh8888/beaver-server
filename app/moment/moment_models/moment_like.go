package moment_models

import (
	"beaver/common/models"
)

/**
 * @description: 点赞表
 */
type MomentLikeModel struct {
	models.Model
	LikeID    string `gorm:"column:like_id;size:64;uniqueIndex;not null" json:"likeId"` // 全局唯一ID
	MomentID  string `gorm:"size:64;not null;index" json:"momentId"`                    // 动态ID (关联 moment_id)
	UserID    string `gorm:"size:64;not null;index" json:"userId"`                      // 点赞用户Id (索引)
	IsDeleted bool   `gorm:"not null;default:false;index" json:"isDeleted"`             // 软删除标记 (索引)
}
