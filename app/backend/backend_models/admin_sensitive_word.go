package backend_models

import "beaver/common/models"

// AdminSensitiveWord 运营敏感词库（后台域）
type AdminSensitiveWord struct {
	models.Model
	Word     string `gorm:"size:128;not null;uniqueIndex" json:"word"`
	Category string `gorm:"size:64;index" json:"category"`
	Level    int    `gorm:"type:tinyint;default:1;index" json:"level"` // 1低 2中 3高
	IsActive bool   `gorm:"default:true;index" json:"isActive"`
	Remark   string `gorm:"size:256" json:"remark"`
}
