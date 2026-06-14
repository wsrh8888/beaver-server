package open_models

import "gorm.io/gorm"

// OpenOAuthToken 用户授权令牌表（含 refresh token）
type OpenOAuthToken struct {
	gorm.Model
	AppID                 string `gorm:"type:varchar(64);index;not null;comment:应用ID"`
	Token                 string `gorm:"type:varchar(256);uniqueIndex;not null;comment:访问令牌"`
	RefreshToken          string `gorm:"type:varchar(256);index;comment:刷新令牌"`
	ExpiresAt             int64  `gorm:"type:bigint;not null;comment:access_token过期时间戳"`
	RefreshTokenExpiresAt int64  `gorm:"type:bigint;not null;comment:refresh_token过期时间戳"`
	Scope                 string `gorm:"type:text;comment:授权范围"`
	UserID                string `gorm:"type:varchar(64);comment:用户ID"`
	OpenID                string `gorm:"type:varchar(64);index;comment:授权用户唯一标识"`
	UnionID               string `gorm:"type:varchar(64);index;comment:用户统一标识"`
}
