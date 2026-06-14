package platform_models

import (
	"beaver/common/models"
)

type UpdateVersion struct {
	models.Model
	ArchitectureID uint   `json:"architectureId" gorm:"index"`
	Version        string `json:"version"`
	FileUrl        string `json:"fileUrl" gorm:"type:varchar(512)"`
	Description    string `json:"description"`
	ReleaseNotes   string `json:"releaseNotes"`
}
