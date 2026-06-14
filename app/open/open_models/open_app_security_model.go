package open_models

import (
	"gorm.io/gorm"
)

// ==================== Security 配置表 ====================

// OpenAppSecurity 安全配置表（对标钉钉开放平台）
type OpenAppSecurity struct {
	gorm.Model
	AppID string `gorm:"type:varchar(64);uniqueIndex;not null;comment:应用ID"`

	// IP 白名单
	IPWhitelist    string `gorm:"type:text;comment:IP白名单(JSON数组)"`
	TrustedDomains string `gorm:"type:text;comment:可信域名列表(JSON数组)"`

	// 限流配置
	RateLimitEnabled bool `gorm:"type:tinyint;default:0;comment:是否启用限流 1是 0否"`
	RateLimitQPS     int  `gorm:"type:int;default:100;comment:每秒请求数限制"`

	// CSRF 保护
	CSRFProtection bool `gorm:"type:tinyint;default:0;comment:是否启用CSRF保护 1是 0否"`

	// CORS 配置
	AllowedOrigins string `gorm:"type:text;comment:CORS允许的源(JSON数组)"`

	// HTTPS 强制
	RequireHTTPS bool `gorm:"type:tinyint;default:0;comment:是否强制HTTPS 1是 0否"`
}
