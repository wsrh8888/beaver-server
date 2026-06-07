package auth_models

import (
	"beaver/common/models"
	"time"
)

// AuthDeviceModel 登录设备档案（MySQL 持久化：设备列表、踢下线、多端登录）。
// 与 WS 实时在线态（Redis coreonline）分离：这里是「登记过哪些设备」，不是「此刻是否连着 WS」。
type AuthDeviceModel struct {
	models.Model
	UserID        string    `gorm:"size:64;not null;uniqueIndex:idx_auth_user_device" json:"userId"`
	DeviceID      string    `gorm:"size:128;not null;uniqueIndex:idx_auth_user_device" json:"deviceId"`
	DeviceType    string    `gorm:"size:32;not null;index" json:"deviceType"` // 槽位：desktop / mobile
	DeviceOS      string    `gorm:"size:32;index" json:"deviceOs"`
	DeviceName    string    `gorm:"size:128" json:"deviceName"`
	DeviceInfo    string    `gorm:"type:text" json:"deviceInfo"`
	LastLoginTime time.Time `json:"lastLoginTime"`
	IsActive      bool      `gorm:"not null;default:true;index" json:"isActive"`
	LastLoginIP   string    `gorm:"size:39" json:"lastLoginIp"`
}
