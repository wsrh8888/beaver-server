package call_models

import (
	"beaver/common/models"
	"time"
)

// CallParticipant 通话参与者表
type CallParticipant struct {
	models.Model
	RoomID string `gorm:"type:varchar(64);index:idx_room_user;not null;comment:关联RoomID" json:"room_id"`
	UserID string `gorm:"type:varchar(64);index:idx_room_user;not null;comment:用户ID" json:"user_id"`
	// 核心行为状态
	Status int8 `gorm:"type:tinyint;default:1;comment:状态:1-进行中,2-已结束" json:"status"`
	Role   int8 `gorm:"type:tinyint;default:1;comment:角色:1-发起者,2-受邀者" json:"role"`

	JoinTime  *time.Time `gorm:"type:datetime;comment:加入时间" json:"join_time"`
	LeaveTime *time.Time `gorm:"type:datetime;comment:离开时间" json:"leave_time"`

	// 扩展信息（用于质量监控，大厂通常会有）
	DeviceInfo string `gorm:"type:varchar(255);comment:设备信息(iOS/Android)" json:"device_info"`
}
