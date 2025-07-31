package moment_models

import (
	"beaver/app/user/user_models"
	"beaver/common/models"
)

/**
 * @description: 点赞表
 */
type MomentLikeModel struct {
	models.Model
	MomentID      uint                  `gorm:"size:64;index;not null" json:"momentId"` // 动态Id
	UserID        string                `gorm:"size:64;index" json:"userId"`            // 接收验证方的 UserID
	LikeUserModel user_models.UserModel `gorm:"foreignKey:UserID;references:UUID" json:"-"`
}
