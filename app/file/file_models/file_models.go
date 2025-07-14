package file_models

import "beaver/common/models"

type FileModel struct {
	models.Model
	FileID   string `json:"fileId"`   // 文件唯一ID /api/file/{uuid}
	FileName string `json:"fileName"` // 文件名
	Size     int64  `json:"size"`     // 文件大小
	Path     string `json:"path"`     // 文件实际存储路径
	Md5      string `json:"md5"`      // 文件md5
	Type     string `json:"type"`     // 文件类型
}
