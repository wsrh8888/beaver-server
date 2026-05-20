package open_models

import (
	"gorm.io/gorm"
)

// OpenEventSubscription 事件订阅配置表
type OpenEventSubscription struct {
	gorm.Model
	AppID      string `gorm:"size:64;index;not null;comment:应用ID"`
	EventType  string `gorm:"size:64;index;not null;comment:事件类型"`
	TargetURL  string `gorm:"size:512;not null;comment:推送目标URL"`
	Secret     string `gorm:"size:128;comment:签名密钥"`
	Status     int    `gorm:"default:1;comment:状态 0禁用 1启用"`
	RetryCount int    `gorm:"default:3;comment:重试次数"`
	Timeout    int    `gorm:"default:5;comment:超时时间(秒)"`
}
