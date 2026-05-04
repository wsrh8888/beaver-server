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
	CreatedAt  int64  `gorm:"not null;comment:创建时间戳"`
	UpdatedAt  int64  `gorm:"not null;comment:更新时间戳"`
}

func (OpenEventSubscription) TableName() string {
	return "open_event_subscriptions"
}

// OpenEventLog 事件推送日志表
type OpenEventLog struct {
	gorm.Model
	EventID      string `gorm:"size:64;uniqueIndex;not null;comment:事件ID"`
	AppID        string `gorm:"size:64;index;not null;comment:应用ID"`
	EventType    string `gorm:"size:64;index;not null;comment:事件类型"`
	Payload      string `gorm:"type:text;comment:事件数据(JSON)"`
	TargetURL    string `gorm:"size:512;comment:推送目标URL"`
	ResponseCode int    `gorm:"comment:响应状态码"`
	ResponseTime int    `gorm:"comment:响应时间(ms)"`
	RetryCount   int    `gorm:"default:0;comment:重试次数"`
	Status       int    `gorm:"default:0;comment:状态 0待推送 1成功 2失败"`
	ErrorMsg     string `gorm:"size:512;comment:错误信息"`
	CreatedAt    int64  `gorm:"not null;comment:创建时间戳"`
}

func (OpenEventLog) TableName() string {
	return "open_event_logs"
}
