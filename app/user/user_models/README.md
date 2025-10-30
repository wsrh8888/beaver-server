1、建议增加 UserLoginHistory 用户登录历史表

type UserLoginHistory struct {
    models.Model
    UserID      string    `gorm:"size:64;not;index" json:"userId"`      // 用户ID
    DeviceID    string    `gorm:"size:128;not" json:"deviceId"`         // 设备唯一标识符
    DeviceType  string    `gorm:"size:32;not;index" json:"deviceType"`  // desktop/mobile/web
    DeviceOS    string    `gorm:"size:32;not;index" json:"deviceOs"`    // windows/macos/linux/android/ios/harmonyos
    DeviceName  string    `gorm:"size:128" json:"deviceName"`                // 设备名称
    LoginTime   time.Time `gorm:"not;index" json:"loginTime"`           // 登录时间
    LoginIP     string    `gorm:"size:39" json:"loginIp"`                    // 登录IP
    Location    string    `gorm:"size:128" json:"location"`                  // 登录位置（可选）
    LoginResult int8      `gorm:"default:1" json:"loginResult"`              // 1: 成功 2: 失败
    // 可选扩展：失败原因、认证方式（密码/扫码/OAuth）、是否是异常登录等
}