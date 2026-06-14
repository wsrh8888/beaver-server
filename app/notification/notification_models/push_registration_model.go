package notification_models

import (
	"beaver/common/models"
)

// PushRegistrationModel 离线推送 Token 注册表（notification 库）
type PushRegistrationModel struct {
	models.Model
	UserID       string `gorm:"size:64;not null;uniqueIndex:idx_push_user_device" json:"userId"`
	DeviceID     string `gorm:"size:128;not null;uniqueIndex:idx_push_user_device" json:"deviceId"`
	PushToken    string `gorm:"size:512;not null" json:"-"`
	PushPlatform string `gorm:"size:16;not null;index" json:"pushPlatform"` // fcm | apns
	Enabled      bool   `gorm:"not null;default:true;index" json:"enabled"`
}
