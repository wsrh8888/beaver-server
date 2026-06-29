package circle_models

import "beaver/common/models"

// CircleMemberModel 圈子成员表
type CircleMemberModel struct {
	models.Model
	CircleID string `gorm:"size:64;not null;uniqueIndex:idx_circle_member" json:"circleId"` // 圈子ID
	UserID   string `gorm:"size:64;not null;uniqueIndex:idx_circle_member;index" json:"userId"` // 用户ID
	Role     int8   `gorm:"not null;default:3" json:"role"`                                     // 角色：1=圈主 2=管理员 3=普通成员
	Version  int64  `gorm:"not null;default:0;index" json:"version"`                            // 版本号，用于客户端增量同步
}
