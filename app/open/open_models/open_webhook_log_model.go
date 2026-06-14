package open_models

import "gorm.io/gorm"

// OpenWebhookLog Webhook 推送日志
type OpenWebhookLog struct {
	gorm.Model
	SubscriptionID uint   `gorm:"index;comment:订阅ID"`
	ConfigID       string `gorm:"type:varchar(64);index;comment:配置ID(兼容)"`
	AppID          string `gorm:"type:varchar(64);index;comment:应用ID"`
	EventID        string `gorm:"type:varchar(64);index;comment:事件ID"`
	EventType      string `gorm:"type:varchar(100);comment:事件类型"`
	TargetURL      string `gorm:"type:varchar(512);comment:目标URL"`
	HTTPStatus     int    `gorm:"comment:HTTP状态码"`
	LatencyMs      int64  `gorm:"comment:耗时毫秒"`
	RetryCount     int    `gorm:"comment:重试次数"`
	ErrorMessage   string `gorm:"type:varchar(512);comment:错误信息"`
	Status         int    `gorm:"type:tinyint;default:0;comment:1成功 0失败"`
}

func (OpenWebhookLog) TableName() string {
	return "open_webhook_logs"
}
