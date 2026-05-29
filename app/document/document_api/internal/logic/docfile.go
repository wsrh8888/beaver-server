package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/document/document_models"
	"beaver/app/file/file_rpc/types/file_rpc"
	"gorm.io/gorm"
)

func personalSpaceID(userID string) string {
	return fmt.Sprintf("user:%s", userID)
}

func ensureUserSpace(db *gorm.DB, userID string) (string, error) {
	spaceID := personalSpaceID(userID)
	var space document_models.CloudDocumentSpace
	err := db.Where("space_id = ?", spaceID).First(&space).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		space = document_models.CloudDocumentSpace{
			SpaceID:     spaceID,
			SpaceType:   1,
			OwnerID:     userID,
			CreatorID:   userID,
			Name:        "我的云文档",
			DefaultPerm: document_models.DocPermEdit,
			Status:      1,
		}
		if err := db.Create(&space).Error; err != nil {
			return "", err
		}
		return spaceID, nil
	}
	if err != nil {
		return "", err
	}
	return spaceID, nil
}

func validateDocumentFile(ctx context.Context, fileRpc file_rpc.FileClient, fileKey string) (*file_rpc.GetFileDetailRes, error) {
	if fileKey == "" {
		return nil, fmt.Errorf("文件标识不能为空")
	}
	detail, err := fileRpc.GetFileDetail(ctx, &file_rpc.GetFileDetailReq{FileKey: fileKey})
	if err != nil {
		return nil, fmt.Errorf("正文文件不存在")
	}
	return detail, nil
}
