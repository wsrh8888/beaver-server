package circle_models

import "beaver/common/models"

// CircleLikeModel 圈子帖子点赞表
type CircleLikeModel struct {
	models.Model
	PostID   string `gorm:"size:64;not null;uniqueIndex:idx_post_user_like;index" json:"postId"` // 帖子ID
	UserID   string `gorm:"size:64;not null;uniqueIndex:idx_post_user_like" json:"userId"`       // 点赞用户ID
	CircleID string `gorm:"size:64;not null;index" json:"circleId"`                             // 所属圈子ID（冗余便于查询）
}
