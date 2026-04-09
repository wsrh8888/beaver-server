package friend_models

import "beaver/common/models"

// FriendBlockModel 黑名单表
type FriendBlockModel struct {
	models.Model
	BlockID       string `gorm:"column:block_id;size:64;uniqueIndex;not null" json:"blockId"`
	UserID        string `gorm:"size:64;index;not null" json:"userId"`        // 执行拉黑的用户ID
	BlockedUserID string `gorm:"size:64;index;not null" json:"blockedUserId"` // 被拉黑的用户ID
}
