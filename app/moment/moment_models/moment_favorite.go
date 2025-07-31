package moment_models

import (
	"beaver/common/models"
)

/**
 * @description: 动态收藏表
 */
type MomentFavoriteModel struct {
	models.Model
	UserID   string `gorm:"size:64;not null" json:"userId"` // 收藏用户ID
	MomentID uint   `gorm:"not null" json:"momentId"`       // 被收藏的动态ID
}
