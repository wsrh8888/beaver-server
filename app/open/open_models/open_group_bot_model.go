package open_models

import "gorm.io/gorm"

type OpenGroupBotModel struct {
	gorm.Model
	Token     string `gorm:"type:varchar(128);uniqueIndex;not null;comment:Webhook Token（URL参数）"`
	Secret    string `gorm:"type:varchar(128);not null;comment:HMAC-SHA256 签名密钥"`
	AppID     string `gorm:"type:varchar(64);index;not null;comment:应用ID"`
	GroupID   string `gorm:"type:varchar(64);index;not null;comment:目标群组ID"`
	BotUserID string `gorm:"type:varchar(64);not null;comment:Bot的UserID"`
	Status    int    `gorm:"type:tinyint;default:1;comment:状态 1启用 0禁用"`
}
