package update_models

import (
	"beaver/common/models"
)

// 架构信息表 - 用于存储不同平台的架构信息
type UpdateArchitecture struct {
	models.Model
	AppID       string          `json:"appId" gorm:"type:varchar(64);index"`         // 关联的应用ID
	App         *UpdateApp      `json:"app" gorm:"foreignKey:AppID;references:UUID"` // 关联的应用信息
	PlatformID  uint            `json:"platformId"`                                  // 平台ID：1=Windows, 2=MacOS, 3=iOS, 4=Android, 5=HarmonyOS
	ArchID      uint            `json:"archId"`                                      // 架构类型 1=WinX64, 2=WinArm64, 3=MacIntel, 4=MacApple, 5=iOS, 6=Android, 7=HarmonyOS
	Description string          `json:"description"`                                 // 架构描述
	Versions    []UpdateVersion `gorm:"foreignKey:ArchitectureID"`                   // 关联的版本信息
	IsActive    bool            `json:"isActive"`                                    // 是否活跃
}

// 预定义平台类型
const (
	PlatformWindows   uint = 1 // Windows平台
	PlatformMacOS     uint = 2 // MacOS平台（苹果电脑）
	PlatformIOS       uint = 3 // iOS平台（苹果手机/平板）
	PlatformAndroid   uint = 4 // Android平台
	PlatformHarmonyOS uint = 5 // 鸿蒙系统平台
)

// 预定义架构类型常量 - 按平台区分
const (
	// H5架构
	H5 uint = 0 // H5网页版本

	// Windows架构
	WinX64   uint = 1 // Intel x64 (64位)
	WinArm64 uint = 2 // ARM64 (Surface等设备)

	// MacOS架构
	MacIntel uint = 3 // Intel版本
	MacApple uint = 4 // Apple Silicon (M1/M2/M3系列)

	// iOS架构
	IOS uint = 5 // iOS通用版本

	// Android架构
	Android uint = 6 // 通用版本

	// 鸿蒙系统架构
	HarmonyOS uint = 7 // 鸿蒙系统通用版本
)
