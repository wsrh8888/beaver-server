package update_models

import "beaver/common/models"

// 应用信息表 - 顶层
type UpdateApp struct {
	models.Model
	Name        string `json:"name" gorm:"size:64"` // 应用名称，如"飞书"、"微信"
	Description string `json:"description"`         // 应用描述
	UUID        string `gorm:"size:64;unique;index"`
	Icon        string `json:"icon"`     // 应用图标URL
	IsActive    bool   `json:"isActive"` // 是否活跃
}
