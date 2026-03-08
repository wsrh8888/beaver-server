package call_models

import (
	"beaver/common/models"
	"time"
)

/*
ParticipantStatus 业务逻辑说明：
1. 待接听 (1): 正在被呼叫，手机正在振铃中。
2. 已接听 (2): 用户点击接听，成功进入音视频房间。
3. 拒绝 (3): 振铃阶段用户手动点击“挂断”或“拒绝”。
4. 超时未接 (4): 振铃超时，用户未操作，系统自动触发取消。
5. 已退出 (5): 曾经进场，后续正常离开、手动挂断或掉线。
*/

// ParticipantStatus 参与者状态
type ParticipantStatus int8

const (
	ParticipantStatusCalling  ParticipantStatus = 1 // 待接听
	ParticipantStatusJoined   ParticipantStatus = 2 // 已接听
	ParticipantStatusRejected ParticipantStatus = 3 // 拒绝
	ParticipantStatusTimeout  ParticipantStatus = 4 // 超时未接
	ParticipantStatusLeft     ParticipantStatus = 5 // 已退出
)

// CallParticipant 通话参与者表
type CallParticipant struct {
	models.Model
	RoomID string `gorm:"type:varchar(64);index:idx_room_user;not null;comment:关联RoomID" json:"room_id"`
	UserID string `gorm:"type:varchar(64);index:idx_room_user;not null;comment:用户ID" json:"user_id"`
	// 核心行为状态
	Status ParticipantStatus `gorm:"type:tinyint;default:1;comment:状态:1-待接听,2-已接听,3-拒绝,4-超时,5-挂断" json:"status"`
	Role   int8              `gorm:"type:tinyint;default:1;comment:角色:1-发起者,2-受邀者" json:"role"`

	JoinTime  *time.Time `gorm:"type:datetime;comment:加入时间" json:"join_time"`
	LeaveTime *time.Time `gorm:"type:datetime;comment:离开时间" json:"leave_time"`

	// 扩展信息（用于质量监控，大厂通常会有）
	DeviceInfo string `gorm:"type:varchar(255);comment:设备信息(iOS/Android)" json:"device_info"`
}
