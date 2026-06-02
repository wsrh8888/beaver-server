package platform_models

import "beaver/common/models"

// UpdateApp 应用信息表
type UpdateApp struct {
	models.Model
	Name        string `json:"name" gorm:"size:64"`
	Description string `json:"description"`
	AppID       string `gorm:"column:app_id;size:64;uniqueIndex" json:"appId"`
	Icon        string `json:"icon"`
	IsActive    bool   `json:"isActive"`
}
