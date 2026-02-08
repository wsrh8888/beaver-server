package call_models

import (
	"beaver/common/models"
	"time"
)

// CallSession 通话主表（对应一次呼叫流程）
type CallSession struct {
	models.Model
	RoomID         string     `gorm:"type:varchar(64);uniqueIndex;not null;comment:LiveKit房间ID/业务唯一ID" json:"room_id"`
	CallerID       string     `gorm:"type:varchar(64);index;not null;comment:发起者ID" json:"caller_id"`
	ConversationID string     `gorm:"type:varchar(64);index;comment:会话ID" json:"conversation_id"`
	CallType       int8       `gorm:"type:tinyint;default:1;comment:通话类型:1-私聊,2-群聊" json:"call_type"`
	Status         int8       `gorm:"type:tinyint;default:1;index;comment:状态:1-进行中,2-已结束" json:"status"`
	StartTime      *time.Time `gorm:"type:datetime;comment:接通时间" json:"start_time"`
	EndTime        *time.Time `gorm:"type:datetime;comment:挂断时间" json:"end_time"`
	Duration       int32      `gorm:"type:int;default:0;comment:通话时长(秒)" json:"duration"`
}
