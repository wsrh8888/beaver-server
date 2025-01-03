package file_models

import "beaver/common/models"

type FileModel struct {
	models.Model
	FileId   string `json:"fileId"`   // 文件唯一Id /api/file/{uuid}
	UserId   string `json:"userId"`   // 用户Id
	FileName string `json:"fileName"` // 文件名
	Size     int64  `json:"size"`     // 文件大小
	Path     string `json:"path"`     // 文件实际存储路径
	Hash     string `json:"hash"`     // 文件hash
	Type     string `json:"type"`     // 文件类型
}
