package open_models

import "beaver/common/models"

// OpenOAuthCode OAuth 授权码表（Scene: oauth/h5_sso/pc_scan）
type OpenOAuthCode struct {
	models.Model
	Code        string `gorm:"size:128;uniqueIndex;not null;comment:授权码"`
	AppID       string `gorm:"size:64;index;not null;comment:应用ID"`
	UserID      string `gorm:"size:64;not null;comment:用户ID"`
	RedirectURI string `gorm:"size:256;comment:回调地址"`
	Scope       string `gorm:"size:256;comment:权限范围"`
	State       string `gorm:"size:128;comment:CSRF state"`
	ExpiresAt   int64  `gorm:"not null;comment:过期时间"`
	Used        bool   `gorm:"default:false;comment:是否已使用"`
	OpenID      string `gorm:"size:64;comment:授权用户唯一标识"`
	UnionID     string `gorm:"size:64;comment:用户统一标识"`
	Scene       string `gorm:"size:20;default:oauth;comment:场景 oauth/h5_sso/pc_scan"`
	IP          string `gorm:"size:64;comment:客户端IP"`
	UserAgent   string `gorm:"size:512;comment:User-Agent"`
}
