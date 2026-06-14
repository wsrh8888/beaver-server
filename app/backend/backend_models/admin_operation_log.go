package backend_models

import "beaver/common/models"

// AdminOperationLog 管理员操作审计日志
type AdminOperationLog struct {
	models.Model
	OperatorID   string `gorm:"size:64;index;not null" json:"operatorId"`
	Action       string `gorm:"size:64;index;not null" json:"action"`
	TargetType   string `gorm:"size:32;index" json:"targetType"`
	TargetID     string `gorm:"size:128;index" json:"targetId"`
	CaseID       uint64 `gorm:"default:0;index" json:"caseId"`
	Detail       string `gorm:"type:text" json:"detail"`
	Result       string `gorm:"size:32;not null" json:"result"`
	ErrorMessage string `gorm:"type:text" json:"errorMessage"`
}
