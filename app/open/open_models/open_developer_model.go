package open_models

import "beaver/common/models"

// OpenDeveloper 开发者申请表
type OpenDeveloper struct {
	models.Model
	UserID      string `gorm:"size:64;uniqueIndex;not null;comment:用户ID"`
	RealName    string `gorm:"size:32;comment:真实姓名"`
	CompanyName string `gorm:"size:64;comment:公司名称"`
	Phone       string `gorm:"size:11;comment:联系电话"`
	Email       string `gorm:"size:128;comment:邮箱"`
	Description string `gorm:"type:text;comment:申请说明"`
	Status      int    `gorm:"default:0;comment:状态 0待审核 1已通过 2已拒绝"`
	AuditBy     string `gorm:"size:64;comment:审核人ID"`
	AuditTime   int64  `gorm:"comment:审核时间"`
	AuditRemark string `gorm:"type:text;comment:审核备注"`
}
