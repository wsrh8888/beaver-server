package open_models

import (
	"gorm.io/gorm"
)

// OpenOAuthConfig OAuth 配置表
type OpenOAuthConfig struct {
	gorm.Model
	AppID           string `gorm:"type:varchar(64);uniqueIndex;not null;comment:应用ID"`
	RedirectURIs    string `gorm:"type:text;comment:回调地址列表(JSON数组)"`
	Scopes          string `gorm:"type:text;comment:授权范围(JSON数组)"`
	CustomLogo      string `gorm:"type:varchar(500);comment:自定义Logo URL"`
	CustomTitle     string `gorm:"type:varchar(100);comment:自定义标题"`
	CustomColor     string `gorm:"type:varchar(20);comment:主题颜色"`
	EnablePKCE      int    `gorm:"type:tinyint;default:0;comment:是否启用PKCE 1是 0否"`
	TokenExpiration int    `gorm:"type:int;default:7200;comment:Token过期时间(秒)"`
	Status          int    `gorm:"type:tinyint;default:1;comment:状态 1启用 0禁用"`
}
