package update_models

import (
	"beaver/common/models"
	"time"
)

// 平台模型
type UpdatePlatform struct {
	models.Model
	Name                  string          `json:"name"`                  // 平台名称，如 iOS, Android, Windows, MacOS, Linux
	Description           string          `json:"description"`           // 平台描述
	Versions              []UpdateVersion `gorm:"foreignkey:PlatformID"` // 关联的版本信息
	MinimumSupportVersion string          `json:"minimumSupportVersion"` // 最低支持版本
	Architecture          []string        `json:"architecture"`          // 支持的架构，如 arm64, x86_64
	PackageType           string          `json:"packageType"`           // 安装包类型，如 apk, ipa, exe, dmg
	IsActive              bool            `json:"isActive"`              // 该平台是否仍在维护
	UpdateChannel         string          `json:"updateChannel"`         // 更新渠道，如 appstore, googleplay, direct
}

// 版本模型
type UpdateVersion struct {
	models.Model
	PlatformID      uint      `json:"platformId"`      // 关联的平台ID
	Version         string    `json:"version"`         // 版本号
	BuildNumber     int       `json:"buildNumber"`     // 构建号
	UpdateType      string    `json:"updateType"`      // 更新类型：optional（可选更新）, mandatory（强制更新）
	DownloadURL     string    `json:"downloadUrl"`     // 下载链接
	Description     string    `json:"description"`     // 版本描述
	ReleaseNotes    string    `json:"releaseNotes"`    // 更新日志
	ReleaseDate     time.Time `json:"releaseDate"`     // 发布时间
	Size            int64     `json:"size"`            // 安装包大小（字节）
	MD5             string    `json:"md5"`             // 安装包MD5校验
	Signature       string    `json:"signature"`       // 安装包签名
	CompatibleRange string    `json:"compatibleRange"` // 兼容版本范围
	IsGrayRelease   bool      `json:"isGrayRelease"`   // 是否灰度发布
	GrayScale       int       `json:"grayScale"`       // 灰度比例（0-100）
}

// 用户版本上报模型
type UpdateUserVersionReport struct {
	models.Model
	UserID         uint      `json:"userId"`         // 用户ID
	Platform       string    `json:"platform"`       // 平台名称
	Version        string    `json:"version"`        // 用户版本号
	BuildNumber    int       `json:"buildNumber"`    // 构建号
	DeviceInfo     string    `json:"deviceInfo"`     // 设备信息
	SystemVersion  string    `json:"systemVersion"`  // 系统版本
	Architecture   string    `json:"architecture"`   // CPU架构
	UpdateChannel  string    `json:"updateChannel"`  // 更新渠道
	LastCheckTime  time.Time `json:"lastCheckTime"`  // 最后检查更新时间
	LastUpdateTime time.Time `json:"lastUpdateTime"` // 最后更新时间
	UpdateStatus   string    `json:"updateStatus"`   // 更新状态：up_to_date, pending, failed
}
