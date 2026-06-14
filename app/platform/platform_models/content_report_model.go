package platform_models

import (
	"beaver/common/models"
)

// ContentReportTargetType 举报对象类型
const (
	ReportTargetUser    = 1
	ReportTargetMessage = 2
	ReportTargetMoment  = 3
	ReportTargetGroup   = 4
)

// ContentReportReason 举报原因
const (
	ReportReasonSpam       = 1
	ReportReasonHarassment = 2
	ReportReasonIllegal    = 3
	ReportReasonOther      = 4
)

// ContentReportStatus 举报状态
const (
	ReportStatusPending  = 1
	ReportStatusAccepted = 2 // 已立案（关联工单）
	ReportStatusRejected = 3
	ReportStatusResolved = 4
)

// ContentReportModel 内容举报表（C 端提交，运营后台处置）
type ContentReportModel struct {
	models.Model
	ReporterUserID string    `gorm:"size:64;index;not null" json:"reporterUserId"`
	TargetType     int       `gorm:"type:tinyint;not null;index" json:"targetType"`
	TargetID       string    `gorm:"size:128;not null;index" json:"targetId"`
	ReasonType     int       `gorm:"type:tinyint;not null" json:"reasonType"`
	Content        string    `gorm:"type:text" json:"content"`
	FileNames      FileNames `gorm:"type:json" json:"fileNames"`
	Status         int       `gorm:"type:tinyint;not null;default:1;index" json:"status"`
	CaseID         uint64    `gorm:"default:0;index" json:"caseId"`
	HandlerID      string    `gorm:"size:64" json:"handlerId"`
	HandleRemark   string    `gorm:"type:text" json:"handleRemark"`
}
