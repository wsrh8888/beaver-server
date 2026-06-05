package logic

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"beaver/app/file/file_models"
	"beaver/app/file/file_rpc/internal/svc"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type SaveFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveFileLogic {
	return &SaveFileLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SaveFileLogic) SaveFile(in *file_rpc.SaveFileReq) (*file_rpc.SaveFileRes, error) {
	var existing file_models.FileModel
	if err := l.svcCtx.DB.Take(&existing, "md5 = ?", in.Md5).Error; err == nil {
		return &file_rpc.SaveFileRes{FileKey: existing.FileKey}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if in.FileInfoJson == "" {
		return nil, status.Error(codes.InvalidArgument, "fileInfo不能为空")
	}

	suffix := "jpg"
	if strings.Contains(in.OriginalName, ".") {
		parts := strings.Split(in.OriginalName, ".")
		if len(parts) > 1 {
			suffix = strings.ToLower(parts[len(parts)-1])
		}
	}
	fileKey := in.Md5 + "." + suffix

	source := file_models.QiniuSource
	if in.Source == "local" {
		source = file_models.LocalSource
	}

	fileInfo := &file_models.FileInfo{}
	if err := json.Unmarshal([]byte(in.FileInfoJson), fileInfo); err != nil {
		return nil, status.Error(codes.InvalidArgument, "fileInfo格式不正确")
	}

	newFile := &file_models.FileModel{
		FileKey:      fileKey,
		OriginalName: in.OriginalName,
		Size:         in.Size,
		Path:         in.Path,
		Md5:          in.Md5,
		Type:         in.Type,
		Source:       source,
		FileInfo:     fileInfo,
	}
	if err := l.svcCtx.DB.Create(newFile).Error; err != nil {
		l.Errorf("保存文件失败: %v", err)
		return nil, status.Error(codes.Internal, "保存文件失败")
	}

	return &file_rpc.SaveFileRes{FileKey: fileKey}, nil
}
