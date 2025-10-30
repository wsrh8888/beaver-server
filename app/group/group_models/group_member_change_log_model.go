package group_models

import (
	"beaver/common/models"
	"time"
)

// 群成员变更日志
type GroupMemberChangeLogModel struct {
	models.Model
	GroupID    string    `gorm:"size:64;index" json:"groupId"` // 群ID
	UserID     string    `gorm:"size:64;index" json:"userId"`  // 用户ID
	ChangeType string    `gorm:"size:32" json:"changeType"`    // join/leave/kick/promote/demote
	OperatedBy string    `gorm:"size:64" json:"operatedBy"`    // 操作者（群主/管理员）
	ChangeTime time.Time `json:"changeTime"`                   // 变更时间
	Version    int64     `gorm:"not;index" json:"version"`     // 版本号
}
