package user_models

import (
	"beaver/common/models"
)

// UserChangeLogModel 用户变更日志模型
type UserChangeLogModel struct {
	models.Model
	UserID     string `gorm:"size:64;not;index" json:"userId"`     // 变更的用户ID
	ChangeType string `gorm:"size:32;not;index" json:"changeType"` // 变更类型：nickname/avatar/abstract/gender/status
	NewValue   string `gorm:"type:text" json:"newValue"`           // 变更后的值
	ChangeTime int64  `gorm:"not;index" json:"changeTime"`         // 变更时间戳
	Version    int64  `gorm:"not;index" json:"version"`            // 用户独立递增（与UserModel.Version一样）
}
