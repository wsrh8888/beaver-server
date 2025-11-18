package update_models

import (
	"beaver/common/models"
)

type UpdateVersion struct {
	models.Model
	ArchitectureID uint   `json:"architectureId" gorm:"index"` // 架构id (统一小写d)
	Version        string `json:"version"`                     // 版本号
	FileKey        string `json:"fileKey"`                     // 文件Key
	Description    string `json:"description"`                 // 版本描述
	ReleaseNotes   string `json:"releaseNotes"`                // 更新日志

	// 关联关系
	Architecture *UpdateArchitecture `gorm:"foreignKey:ArchitectureID"` // 架构关联
}
