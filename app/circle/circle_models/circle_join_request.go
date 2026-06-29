package circle_models

import "beaver/common/models"

// CircleJoinRequestModel 加圈申请表
type CircleJoinRequestModel struct {
	models.Model
	CircleID string `gorm:"size:64;not null;index" json:"circleId"`                           // 圈子ID
	UserID   string `gorm:"size:64;not null;index" json:"userId"`                             // 申请用户ID
	Status   int8   `gorm:"not null;default:0" json:"status"`                                 // 状态：0=待审批 1=已通过 2=已拒绝
	Reason   string `gorm:"size:256" json:"reason"`                                           // 申请理由
}
