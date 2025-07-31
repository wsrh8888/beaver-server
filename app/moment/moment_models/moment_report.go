package moment_models

import (
	"beaver/common/models"
)

/**
 * @description: 动态举报表
 */
type MomentReportModel struct {
	models.Model
	UserID   string `gorm:"size:64;not null" json:"userId"`   // 举报用户ID
	MomentID uint   `gorm:"not null" json:"momentId"`         // 被举报的动态ID
	Reason   string `gorm:"type:text;not null" json:"reason"` // 举报原因
	Images   *Files `gorm:"type:longtext" json:"images"`      // 举报图片
	Status   int    `gorm:"not null;default:0" json:"status"` // 处理状态：0-待处理 1-已处理 2-已驳回
}
