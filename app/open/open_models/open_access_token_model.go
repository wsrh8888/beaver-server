package open_models

import (
	"gorm.io/gorm"
)

// OpenAccessToken 访问令牌表
type OpenAccessToken struct {
	gorm.Model
	AppID        string `gorm:"type:varchar(64);index;not null;comment:应用ID"`
	Token        string `gorm:"type:varchar(256);uniqueIndex;not null;comment:访问令牌"`
	RefreshToken string `gorm:"type:varchar(256);index;comment:刷新令牌"`
	ExpiresAt    int64  `gorm:"type:bigint;not null;comment:过期时间戳"`
	Scope        string `gorm:"type:text;comment:授权范围"`
	UserID       string `gorm:"type:varchar(64);comment:用户ID(用户级token才有)"`
}
