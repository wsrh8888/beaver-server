package platform_models

import (
	"beaver/common/models"
)

// UpdateArchitecture 架构信息表
type UpdateArchitecture struct {
	models.Model
	AppID       string          `json:"appId" gorm:"type:varchar(64);index"`
	App         *UpdateApp      `json:"app" gorm:"foreignKey:AppID;references:AppID"`
	PlatformID  uint            `json:"platformId"`
	ArchID      uint            `json:"archId"`
	Description string          `json:"description"`
	Versions    []UpdateVersion `gorm:"foreignKey:ArchitectureID"`
	IsActive    bool            `json:"isActive"`
}

const (
	PlatformWindows   uint = 1
	PlatformMacOS     uint = 2
	PlatformIOS       uint = 3
	PlatformAndroid   uint = 4
	PlatformHarmonyOS uint = 5
)

const (
	H5        uint = 0
	WinX64    uint = 1
	WinArm64  uint = 2
	MacIntel  uint = 3
	MacApple  uint = 4
	IOS       uint = 5
	Android   uint = 6
	HarmonyOS uint = 7
)
