package group_models

import (
	"beaver/common/models"
	"time"
)

// 按照群组独立递增版本
// GroupJoinRequestModel 入群申请表
type GroupJoinRequestModel struct {
	models.Model
	GroupID         string     `gorm:"size:64;index" json:"groupId"`
	ApplicantUserID string     `gorm:"size:64;index" json:"applicantUserId"`
	Message         string     `gorm:"type:text" json:"message"`
	Status          int8       `gorm:"not null;default:0" json:"status"` // 0待审 1同意 2拒绝
	HandledBy       string     `gorm:"size:64" json:"handledBy"`
	HandledAt       *time.Time `json:"handledAt"` // 处理时间
	Version         int64      `gorm:"not null;default:0;index" json:"version"`
}
