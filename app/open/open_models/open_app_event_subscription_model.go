package open_models

import (
	"gorm.io/gorm"
)

// ==================== Event Subscription 配置表 ====================

// OpenAppEventSubscription 应用事件订阅表（对标飞书开放平台）
type OpenAppEventSubscription struct {
	gorm.Model
	AppID       string `gorm:"type:varchar(64);index;not null;comment:应用ID"`
	EventType   string `gorm:"type:varchar(100);not null;comment:事件类型 im.message.receive_v1"`
	CallbackURL string `gorm:"type:varchar(512);not null;comment:回调地址(Webhook URL)"`
	Secret      string `gorm:"type:varchar(128);comment:签名密钥"`
	Status      int    `gorm:"type:tinyint;default:1;comment:状态 1启用 0禁用"`
	RetryCount  int    `gorm:"type:int;default:3;comment:重试次数"`
	Timeout     int    `gorm:"type:int;default:5;comment:超时时间(秒)"`
}
