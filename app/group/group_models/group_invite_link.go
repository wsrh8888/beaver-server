package group_models

import "beaver/common/models"

// GroupInviteLinkModel 群分享邀请凭证（与成员邀请 GroupInvite 区分）
type GroupInviteLinkModel struct {
	models.Model
	Token     string `gorm:"column:token;size:64;uniqueIndex;not null" json:"token"`
	GroupID   string `gorm:"column:group_id;size:64;not null;index" json:"groupId"`
	CreatorID string `gorm:"column:creator_id;size:64;not null;index" json:"creatorId"`
	ExpireAt  int64  `gorm:"column:expire_at;not null;index" json:"expireAt"`
	MaxUses   int64  `gorm:"column:max_uses;not null;default:0" json:"maxUses"`
	UsedCount int64  `gorm:"column:used_count;not null;default:0" json:"usedCount"`
	Status    int8   `gorm:"column:status;not null;default:1;index" json:"status"` // 1有效 2吊销 3用尽
}
