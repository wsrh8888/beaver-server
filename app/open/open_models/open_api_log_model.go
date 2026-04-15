package open_models

import (
	"gorm.io/gorm"
)

// OpenAPILog API 调用日志表
type OpenAPILog struct {
	gorm.Model
	AppID        string `gorm:"type:varchar(64);index;comment:应用ID"`
	APIPath      string `gorm:"type:varchar(200);index;comment:API路径"`
	Method       string `gorm:"type:varchar(10);comment:HTTP方法"`
	RequestIP    string `gorm:"type:varchar(50);comment:请求IP"`
	ResponseCode int    `gorm:"type:int;comment:响应码"`
	ResponseTime int64  `gorm:"type:bigint;comment:响应时间(ms)"`
	ErrorMessage string `gorm:"type:text;comment:错误信息"`
}
