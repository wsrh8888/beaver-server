package moment_models

import (
	"beaver/app/user/user_models"
	"beaver/common/models"
)

/**
 * @description: 评论表
 */
type MomentCommentModel struct {
	models.Model
	MomentID         uint                  `gorm:"size:64;index;not null" json:"momentId"` // 动态Id
	UserID           string                `gorm:"size:64;index" json:"userId"`            // 评论用户Id
	Content          string                `gorm:"type:text;not null" json:"content"`      // 评论内容
	CommentUserModel user_models.UserModel `gorm:"foreignKey:UserID;references:UUID" json:"-"`
}
