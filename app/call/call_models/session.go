package call_models

import (
	"beaver/common/models"
	"time"
)

// CallSession 通话主表（对应一次呼叫流程）
type CallSession struct {
	models.Model
	RoomID    string     `gorm:"type:varchar(64);uniqueIndex;not null;comment:LiveKit房间ID/业务唯一ID" json:"room_id"`
	CallerID  string     `gorm:"type:varchar(64);index;not null;comment:发起者ID" json:"caller_id"`
	GroupID   string     `gorm:"type:varchar(64);index;comment:群组ID(如果是群聊通话)" json:"group_id"`
	CallType  int8       `gorm:"type:tinyint;default:1;comment:通话类型:1-单聊音频,2-单聊视频,3-群聊音频,4-群聊视频" json:"call_type"`
	Status    int8       `gorm:"type:tinyint;default:1;index;comment:状态:1-呼叫中,2-进行中,3-已结束,4-未接听,5-拒接,6-忙线" json:"status"`
	StartTime *time.Time `gorm:"type:datetime;comment:接通时间" json:"start_time"`
	EndTime   *time.Time `gorm:"type:datetime;comment:挂断时间" json:"end_time"`
	Duration  int32      `gorm:"type:int;default:0;comment:通话时长(秒)" json:"duration"`
}
