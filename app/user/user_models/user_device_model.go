package user_models

import (
	"beaver/common/models"
	"time"
)

// UserDeviceModel 用户设备管理模型
type UserDeviceModel struct {
	models.Model
	UserID        string    `gorm:"size:64;not;index" json:"userId"`                           // 用户ID
	DeviceID      string    `gorm:"size:128;not;index:idx_user_device,unique" json:"deviceId"` // 设备ID
	DeviceType    string    `gorm:"size:32;not;index" json:"deviceType"`                       // 设备类型
	DeviceOS      string    `gorm:"size:32;not;index" json:"deviceOs"`                         // 设备操作系统
	DeviceName    string    `gorm:"size:128" json:"deviceName"`                                // 设备名称
	DeviceInfo    string    `gorm:"type:text" json:"deviceInfo"`                               // 设备详细信息
	LastLoginTime time.Time `json:"lastLoginTime"`                                             // 最后登录时间
	IsActive      bool      `gorm:"not;default:true;index" json:"isActive"`                    // 是否活跃
	LastLoginIP   string    `gorm:"size:39" json:"lastLoginIp"`                                // 最后登录IP
}
