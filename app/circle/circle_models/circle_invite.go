package circle_models

import "beaver/common/models"

// CircleInviteModel 圈子分享邀请凭证
type CircleInviteModel struct {
	models.Model
	Token      string `gorm:"column:token;size:64;uniqueIndex;not null" json:"token"` // 对外邀请凭证
	CircleID   string `gorm:"column:circle_id;size:64;not null;index" json:"circleId"`
	CreatorID  string `gorm:"column:creator_id;size:64;not null;index" json:"creatorId"`
	ExpireAt   int64  `gorm:"column:expire_at;not null;index" json:"expireAt"` // Unix 秒
	MaxUses    int64  `gorm:"column:max_uses;not null;default:0" json:"maxUses"` // 0=不限
	UsedCount  int64  `gorm:"column:used_count;not null;default:0" json:"usedCount"`
	Status     int8   `gorm:"column:status;not null;default:1;index" json:"status"` // 1有效 2吊销 3用尽
}
