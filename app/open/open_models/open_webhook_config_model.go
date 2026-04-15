package open_models

import (
	"gorm.io/gorm"
)

// OpenWebhookConfig Webhook 配置表
type OpenWebhookConfig struct {
	gorm.Model
	AppID      string `gorm:"type:varchar(64);index;not null;comment:应用ID"`
	EventType  string `gorm:"type:varchar(50);index;not null;comment:事件类型"`
	TargetURL  string `gorm:"type:varchar(500);not null;comment:回调地址"`
	Secret     string `gorm:"type:varchar(128);comment:签名密钥"`
	Status     int    `gorm:"type:tinyint;default:1;comment:状态 1启用 0禁用"`
	RetryCount int    `gorm:"type:int;default:3;comment:重试次数"`
	Timeout    int    `gorm:"type:int;default:5;comment:超时秒数"`
	StreamMode int    `gorm:"type:tinyint;default:0;comment:是否流式 1是 0否"`
}
