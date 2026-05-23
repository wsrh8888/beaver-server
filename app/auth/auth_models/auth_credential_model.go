package auth_models

import (
	"beaver/common/models"
	"time"
)

// AuthCredentialModel 认证凭证模型（密码、登录记录等敏感信息）
type AuthCredentialModel struct {
	models.Model
	UserID      string     `gorm:"size:64;uniqueIndex;not null" json:"userId"` // 关联 users.user_id
	Password    string     `gorm:"size:128;not null" json:"-"`                 // 密码哈希（不对外暴露）
	Salt        string     `gorm:"size:32" json:"-"`                           // 盐值（可选，如果用 bcrypt 则不需要）
	LastLoginAt *time.Time `json:"lastLoginAt"`                                // 最后登录时间
	LoginCount  int64      `gorm:"default:0" json:"loginCount"`                // 登录次数
}
