package open_models

import "gorm.io/gorm"

// OpenWebhookLog Webhook 回调日志
type OpenWebhookLog struct {
	gorm.Model
	ConfigID  string `gorm:"type:varchar(64);index;comment:配置ID"`
	AppID     string `gorm:"type:varchar(64);index;comment:应用ID"`
	EventType string `gorm:"type:varchar(100);comment:事件类型"`
	Status    int    `gorm:"type:tinyint;default:0;comment:状态 1成功 0失败"`
}

func (OpenWebhookLog) TableName() string {
	return "open_webhook_logs"
}
