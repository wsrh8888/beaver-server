package logic

import (
	"time"

	"beaver/app/file/file_models"
	"beaver/app/file/file_rpc/types/file_rpc"
)

func toFileItem(f file_models.FileModel) *file_rpc.FileItem {
	return &file_rpc.FileItem{
		Id:           uint64(f.Id),
		FileKey:      f.FileKey,
		OriginalName: f.OriginalName,
		Size:         f.Size,
		Path:         f.Path,
		Md5:          f.Md5,
		Type:         f.Type,
		Source:       string(f.Source),
		CreatedAt:    time.Time(f.CreatedAt).Format(time.RFC3339),
		UpdatedAt:    time.Time(f.UpdatedAt).Format(time.RFC3339),
	}
}
