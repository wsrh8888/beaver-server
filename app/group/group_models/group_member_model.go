package group_models

import (
	"beaver/common/models"
	"time"
)

// 群成员（版本号按群组独立递增）
type GroupMemberModel struct {
	models.Model
	GroupID  string    `gorm:"size:64;index" json:"groupId"`
	UserID   string    `gorm:"size:64;index" json:"userId"`
	Role     int8      `json:"role"`                                    // 1群主 2管理员 3普通成员
	Status   int8      `gorm:"default:1" json:"status"`                 // 1正常 2退出 3被踢
	JoinTime time.Time `json:"joinTime"`                                // 加入时间
	Version  int64     `gorm:"not null;default:0;index" json:"version"` // 群组成员列表版本号
}
