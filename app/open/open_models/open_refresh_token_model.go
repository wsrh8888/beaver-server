package open_models

import (
	"gorm.io/gorm"
)

// OpenRefreshToken 刷新令牌表
type OpenRefreshToken struct {
	gorm.Model
	AppID     string `gorm:"type:varchar(64);index;not null;comment:应用ID"`
	Token     string `gorm:"type:varchar(256);uniqueIndex;not null;comment:刷新令牌"`
	ExpiresAt int64  `gorm:"type:bigint;not null;comment:过期时间戳"`
	UserID    string `gorm:"type:varchar(64);comment:用户ID"`
}
