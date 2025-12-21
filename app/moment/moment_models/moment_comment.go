package moment_models

import (
	"beaver/common/models"
)

/**
 * @description: 评论表
 */
type MomentCommentModel struct {
	models.Model
	CommentID        string `gorm:"column:comment_id;size:64;uniqueIndex;not null" json:"commentId"` // 全局唯一ID
	MomentID         string `gorm:"size:64;not null;index" json:"momentId"`                          // 动态ID (关联 moment_id)
	UserID           string `gorm:"size:64;not null;index" json:"userId"`                            // 评论用户Id (索引)
	Content          string `gorm:"type:text;not null" json:"content"`                               // 评论内容
	ParentID         string `gorm:"size:64;index;default:''" json:"parentId"`                        // 父评论ID，空表示一级评论
	ReplyToCommentID string `gorm:"size:64;index;default:''" json:"replyToCommentId"`                // 被回复的评论ID
	IsDeleted        bool   `gorm:"not null;default:false;index" json:"isDeleted"`                   // 软删除标记 (索引)
}
