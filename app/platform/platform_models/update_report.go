package platform_models

import (
	"beaver/common/models"
)

type UpdateReport struct {
	models.Model
	UserID         string `json:"userId" gorm:"size:64;index"`
	DeviceID       string `json:"deviceId" gorm:"size:64;index"`
	AppID          string `json:"appId" gorm:"size:64;index"`
	ArchitectureID uint   `json:"architectureId" gorm:"index"`
	Version        string `json:"version"`
}
