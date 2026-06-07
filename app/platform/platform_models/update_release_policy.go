package platform_models

import "beaver/common/models"

// UpdateReleasePolicy 架构发版策略（正式版 + 比例灰度 + 最低强更）
type UpdateReleasePolicy struct {
	models.Model
	AppID            string `json:"appId" gorm:"type:varchar(64);index"`
	ArchitectureID   uint   `json:"architectureId" gorm:"uniqueIndex"`
	StableVersionID  uint   `json:"stableVersionId"`
	GrayVersionID    uint   `json:"grayVersionId"`
	RolloutPercent   uint   `json:"rolloutPercent"` // 0-100，命中灰度桶的用户使用 GrayVersionID
	MinVersion       string `json:"minVersion" gorm:"type:varchar(32)"`
	ForceUpdate      bool   `json:"forceUpdate"`
	IsActive         bool   `json:"isActive"`
}
