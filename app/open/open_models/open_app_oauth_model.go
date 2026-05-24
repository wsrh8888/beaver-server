package open_models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// ==================== OAuth 配置表 ====================

// OpenAppOAuth OAuth 配置表（对标钉钉开放平台）
type OpenAppOAuth struct {
	gorm.Model
	AppID string `gorm:"type:varchar(64);uniqueIndex;not null;comment:应用ID"`

	// 核心安全配置
	RedirectURIs    string `gorm:"type:text;comment:允许的回调地址(JSON数组)"`
	RequirePKCE     int    `gorm:"type:tinyint;default:0;comment:是否强制要求PKCE 1是 0否"`
	SupportedScopes string `gorm:"type:text;comment:支持的权限范围(JSON数组)"`

	// Token 有效期
	AccessTokenTTL  int `gorm:"type:int;default:7200;comment:Access Token有效期(秒)"`
	RefreshTokenTTL int `gorm:"type:int;default:2592000;comment:Refresh Token有效期(秒)"`

	// 客户端详细配置（JSON）
	H5      *H5OAuth      `gorm:"type:text" json:"h5"`
	Desktop *DesktopOAuth `gorm:"type:text" json:"desktop"`
	Mobile  *MobileOAuth  `gorm:"type:text" json:"mobile"`
}

// H5OAuth H5 应用 OAuth 配置
type H5OAuth struct {
	Enabled      bool     `json:"enabled"`
	RedirectURIs []string `json:"redirectUris"`
	JsSdkDomains []string `json:"jsSdkDomains"`
}

// Value converts the H5OAuth to a JSON-encoded string for database storage
func (c *H5OAuth) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan converts a JSON-encoded string from the database to a H5OAuth
func (c *H5OAuth) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}

// DesktopOAuth 桌面端 OAuth 配置
type DesktopOAuth struct {
	Enabled      bool   `json:"enabled"`
	CustomScheme string `json:"customScheme"`
}

// Value converts the DesktopOAuth to a JSON-encoded string for database storage
func (c *DesktopOAuth) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan converts a JSON-encoded string from the database to a DesktopOAuth
func (c *DesktopOAuth) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}

// MobileOAuth 移动端 OAuth 配置
type MobileOAuth struct {
	Enabled            bool   `json:"enabled"`
	IOSBundleID        string `json:"iosBundleId"`
	AndroidPackageName string `json:"androidPackageName"`
	UniversalLink      string `json:"universalLink"`
	CustomScheme       string `json:"customScheme"`
}

// Value converts the MobileOAuth to a JSON-encoded string for database storage
func (c *MobileOAuth) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan converts a JSON-encoded string from the database to a MobileOAuth
func (c *MobileOAuth) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}
