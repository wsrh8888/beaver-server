package open_models

import (
	"gorm.io/gorm"
)

// OpenWebhookLog Webhook 发送日志表
type OpenWebhookLog struct {
	gorm.Model
	ConfigID     string `gorm:"type:varchar(64);index;comment:配置ID"`
	AppID        string `gorm:"type:varchar(64);index;comment:应用ID"`
	EventType    string `gorm:"type:varchar(50);comment:事件类型"`
	Payload      string `gorm:"type:text;comment:请求体JSON"`
	ResponseCode int    `gorm:"type:int;comment:HTTP状态码"`
	ResponseBody string `gorm:"type:text;comment:响应体"`
	RetryCount   int    `gorm:"type:int;default:0;comment:重试次数"`
	Status       int    `gorm:"type:tinyint;comment:1成功 0失败"`
}
