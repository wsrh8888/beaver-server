package open_models

import (
	"gorm.io/gorm"
)

// OpenApp 开放平台应用表
type OpenApp struct {
	gorm.Model
	AppID       string `gorm:"type:varchar(64);uniqueIndex;not null;comment:应用唯一标识"`
	AppSecret   string `gorm:"type:varchar(128);not null;comment:应用密钥"`
	Name        string `gorm:"type:varchar(100);not null;comment:应用名称"`
	Description string `gorm:"type:text;comment:应用描述"`
	Icon        string `gorm:"type:varchar(500);comment:应用图标URL"`
	OwnerUserID string `gorm:"type:varchar(64);index;comment:所属用户ID"`
	Status      int    `gorm:"type:tinyint;default:0;comment:状态 0草稿 1已发布 2禁用"`
	// 能力开关（对标飞书）
	EnableBot     int `gorm:"type:tinyint;default:0;comment:是否启用机器人能力 1是 0否"`
	EnableOAuth   int `gorm:"type:tinyint;default:0;comment:是否启用OAuth能力 1是 0否"`
	EnableWebhook int `gorm:"type:tinyint;default:0;comment:是否启用Webhook能力 1是 0否"`
	// 其他配置
	WebhookURL  string `gorm:"type:varchar(500);comment:Webhook回调地址"`
	IPWhitelist string `gorm:"type:text;comment:IP白名单(JSON数组)"`
	Scopes      string `gorm:"type:text;comment:权限范围(JSON数组)"`
}
