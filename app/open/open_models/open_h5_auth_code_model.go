package open_models

import (
	"gorm.io/gorm"
)

// OpenH5AuthCode H5 免登授权码表
type OpenH5AuthCode struct {
	gorm.Model
	AppID     string `gorm:"size:64;index;not null;comment:应用ID"`
	Code      string `gorm:"size:128;uniqueIndex;not null;comment:授权码"`
	UserID    string `gorm:"size:64;index;not null;comment:用户ID"`
	ExpiresAt int64  `gorm:"not null;comment:过期时间戳"`
	Used      int    `gorm:"default:0;comment:是否已使用 0未使用 1已使用"`
	UsedAt    int64  `gorm:"comment:使用时间戳"`
	IP        string `gorm:"size:64;comment:客户端IP"`
	UserAgent string `gorm:"size:512;comment:User-Agent"`
}
