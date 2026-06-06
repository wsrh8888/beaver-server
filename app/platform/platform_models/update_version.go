package platform_models

import (
	"beaver/common/models"
)

type UpdateVersion struct {
	models.Model
	ArchitectureID uint   `json:"architectureId" gorm:"index"`
	Version        string `json:"version"`
	FileKey        string `json:"fileKey"`
	Description    string `json:"description"`
	ReleaseNotes   string `json:"releaseNotes"`
}
