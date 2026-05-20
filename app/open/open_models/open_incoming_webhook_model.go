package open_models

import (
	"gorm.io/gorm"
)

// OpenIncomingWebhook 群通知机器人表
// 外部系统（Jenkins/GitHub/Grafana 等）持有 Webhook URL + Secret，
// 通过 HMAC-SHA256 签名向群推送消息，Beaver 只负责转发，不回复。
type OpenIncomingWebhook struct {
	gorm.Model
	Token     string `gorm:"type:varchar(128);uniqueIndex;not null;comment:Webhook Token（URL参数）"`
	Secret    string `gorm:"type:varchar(128);not null;comment:HMAC-SHA256 签名密钥"`
	AppID     string `gorm:"type:varchar(64);index;not null;comment:应用ID"`
	GroupID   string `gorm:"type:varchar(64);index;not null;comment:目标群组ID"`
	BotUserID string `gorm:"type:varchar(64);not null;comment:Bot的UserID（发消息时的发件人身份）"`
	Name          string `gorm:"type:varchar(100);comment:机器人名称（如 Jenkins Bot）"`
	Description   string `gorm:"type:varchar(500);comment:机器人简介"`
	Avatar        string `gorm:"type:varchar(256);comment:机器人头像文件ID"`
	CreatorUserID string `gorm:"type:varchar(64);comment:创建者用户ID"`
	Status        int    `gorm:"type:tinyint;default:1;comment:状态 1启用 0禁用"`
}
