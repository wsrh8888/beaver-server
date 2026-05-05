package open_models

import (
	"gorm.io/gorm"
)

// OpenIncomingWebhook Incoming Webhook Token 表（用于 Jenkins/GitHub 等外部集成）
type OpenIncomingWebhook struct {
	gorm.Model
	Token     string `gorm:"type:varchar(128);uniqueIndex;not null;comment:Webhook Token"`
	AppID     string `gorm:"type:varchar(64);index;not null;comment:应用ID"`
	GroupID   string `gorm:"type:varchar(64);index;not null;comment:群组ID"`
	BotUserID string `gorm:"type:varchar(64);not null;comment:Bot用户ID"`
	Name      string `gorm:"type:varchar(100);comment:Webhook名称"`
	Status    int    `gorm:"type:tinyint;default:1;comment:状态 1启用 0禁用"`
}
