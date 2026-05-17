package open_models

import (
	"gorm.io/gorm"
)

// OpenAppToken 应用级访问令牌表（对标知音楼 /gettoken 接口）
// 用于应用以自身身份调用开放平台接口（如通讯录管理、发送系统消息）
type OpenAppToken struct {
	gorm.Model
	AppID     string `gorm:"type:varchar(64);uniqueIndex;not null;comment:应用ID"`
	Token     string `gorm:"type:varchar(256);not null;comment:应用访问令牌"`
	ExpiresAt int64  `gorm:"type:bigint;not null;comment:过期时间戳"`
}
