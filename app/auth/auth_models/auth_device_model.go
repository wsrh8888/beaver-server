package auth_models

import (
	"beaver/common/models"
	"time"
)

// AuthDeviceModel 登录设备档案（MySQL 持久化：设备列表、踢下线、多端登录）。
// 与 WS 实时在线态（Redis coreonline）分离：这里是「登记过哪些设备」，不是「此刻是否连着 WS」。
type AuthDeviceModel struct {
	models.Model
	UserID          string    `gorm:"size:64;not null;uniqueIndex:idx_auth_user_device" json:"userId"`
	DeviceID        string    `gorm:"size:128;not null;uniqueIndex:idx_auth_user_device" json:"deviceId"`
	DeviceType      string    `gorm:"size:32;not null;index" json:"deviceType"` // 槽位：desktop / mobile
	DeviceOS        string    `gorm:"size:32;index" json:"deviceOs"`            // windows / ios / android / macos / linux
	DeviceModel     string    `gorm:"size:128;index" json:"deviceModel"`        // iPhone17,3、SM-G991B
	DeviceOsVersion string    `gorm:"size:64" json:"deviceOsVersion"`           // 18.2、10.0.19045
	DeviceName      string    `gorm:"size:128" json:"deviceName"`               // iPhone 17 Pro、DESKTOP-HOME
	DeviceInfo      string    `gorm:"type:text" json:"deviceInfo"`
	LastLoginTime   time.Time `json:"lastLoginTime"`
	IsActive        bool      `gorm:"not null;default:true;index" json:"isActive"`
	LastLoginIP     string    `gorm:"size:39" json:"lastLoginIp"`
}
