package circle_models

import "beaver/common/models"

// CircleCommentModel 圈子帖子评论表
type CircleCommentModel struct {
	models.Model
	CommentID        string `gorm:"column:comment_id;size:64;uniqueIndex;not null" json:"commentId"` // 评论唯一ID
	PostID           string `gorm:"size:64;not null;index" json:"postId"`                            // 所属帖子ID
	CircleID         string `gorm:"size:64;not null;index" json:"circleId"`                          // 所属圈子ID（冗余便于查询）
	UserID           string `gorm:"size:64;not null;index" json:"userId"`                            // 评论用户ID
	Content          string `gorm:"type:text;not null" json:"content"`                               // 评论内容
	ParentID         string `gorm:"size:64;index" json:"parentId"`                                   // 父评论ID，空表示一级评论
	ReplyToCommentID string `gorm:"size:64" json:"replyToCommentId"`                                 // 被回复的评论ID
	ReplyToUserID    string `gorm:"size:64" json:"replyToUserId"`                                    // 被回复用户ID
	IsDeleted        bool   `gorm:"not null;default:false;index" json:"isDeleted"`                   // 软删除标记
}
