package open_models

import (
	"gorm.io/gorm"
)

// OpenH5AuthCode H5 免登临时授权码表
type OpenH5AuthCode struct {
	gorm.Model
	Code      string `gorm:"size:64;uniqueIndex;not null;comment:临时授权码"`
	AppID     string `gorm:"size:64;index;not null;comment:应用ID"`
	UserID    string `gorm:"size:64;not null;comment:用户ID"`
	ExpiresAt int64  `gorm:"not null;comment:过期时间戳"`
	CreatedAt int64  `gorm:"not null;comment:创建时间戳"`
}

func (OpenH5AuthCode) TableName() string {
	return "open_h5_auth_codes"
}
