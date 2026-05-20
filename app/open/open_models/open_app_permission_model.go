package open_models

import (
	"gorm.io/gorm"
)

// OpenAppPermission 应用权限表
type OpenAppPermission struct {
	gorm.Model
	AppID     string `gorm:"type:varchar(64);index;not null;comment:应用ID"`
	Scope     string `gorm:"type:varchar(50);not null;comment:权限范围"`
	GrantedAt int64  `gorm:"type:bigint;comment:授权时间戳"`
	GrantedBy string `gorm:"type:varchar(64);comment:授权人"`
	ExpiresAt int64  `gorm:"type:bigint;comment:过期时间戳"`
}
