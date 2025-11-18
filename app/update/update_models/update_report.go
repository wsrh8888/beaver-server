package update_models

import (
	"beaver/common/models"
)

type UpdateReport struct {
	models.Model
	UserID         string `json:"userId" gorm:"size:64;index"`   // 用户ID
	DeviceID       string `json:"deviceId" gorm:"size:64;index"` // 设备ID
	AppID          string `json:"appId" gorm:"size:64;index"`    // 应用ID
	ArchitectureID uint   `json:"architectureId" gorm:"index"`   // 架构ID
	Version        string `json:"version"`                       // 版本号
}
