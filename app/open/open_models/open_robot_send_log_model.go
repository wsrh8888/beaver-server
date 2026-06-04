package open_models

import "gorm.io/gorm"

// OpenRobotSendLog Robot 发消息幂等记录
type OpenRobotSendLog struct {
	gorm.Model
	AppID          string `gorm:"type:varchar(64);uniqueIndex:uk_app_idem;not null;comment:应用ID"`
	IdempotentKey  string `gorm:"type:varchar(128);uniqueIndex:uk_app_idem;not null;comment:幂等键"`
	MessageID      string `gorm:"type:varchar(64);comment:消息ID"`
	ConversationID string `gorm:"type:varchar(128);comment:会话ID"`
}
