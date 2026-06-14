package backend_models

import (
	"beaver/common/models"
	"time"
)

// 工单来源
const (
	CaseSourceReport  = 1
	CaseSourceFeedback = 2
	CaseSourceManual  = 3
)

// 工单状态
const (
	CaseStatusPending    = 1
	CaseStatusProcessing = 2
	CaseStatusResolved   = 3
	CaseStatusRejected   = 4
)

// 处置对象类型
const (
	CaseTargetUser    = 1
	CaseTargetMessage = 2
	CaseTargetMoment  = 3
	CaseTargetGroup   = 4
)

// AdminModerationCase 运营处置工单（后台域，跨 RPC 编排入口）
type AdminModerationCase struct {
	models.Model
	CaseNo       string     `gorm:"size:32;uniqueIndex;not null" json:"caseNo"`
	Source       int        `gorm:"type:tinyint;not null" json:"source"`
	SourceID     uint64     `gorm:"default:0" json:"sourceId"`
	TargetType   int        `gorm:"type:tinyint;not null;index" json:"targetType"`
	TargetID     string     `gorm:"size:128;not null;index" json:"targetId"`
	Title        string     `gorm:"size:200;not null" json:"title"`
	Description  string     `gorm:"type:text" json:"description"`
	Priority     int        `gorm:"type:tinyint;default:1" json:"priority"`
	Status       int        `gorm:"type:tinyint;not null;default:1;index" json:"status"`
	HandlerID    string     `gorm:"size:64;index" json:"handlerId"`
	HandleRemark string     `gorm:"type:text" json:"handleRemark"`
	HandleTime   *time.Time `json:"handleTime"`
	ActionsTaken string     `gorm:"type:text" json:"actionsTaken"`
}
