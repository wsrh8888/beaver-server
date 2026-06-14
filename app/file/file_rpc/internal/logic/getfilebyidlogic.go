package logic

import (
	"context"
	"errors"

	"beaver/app/file/file_models"
	"beaver/app/file/file_rpc/internal/svc"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GetFileByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileByIdLogic {
	return &GetFileByIdLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetFileByIdLogic) GetFileById(in *file_rpc.GetFileByIdReq) (*file_rpc.GetFileByIdRes, error) {
	var file file_models.FileModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "文件不存在")
		}
		return nil, err
	}
	return &file_rpc.GetFileByIdRes{File: toFileItem(file)}, nil
}
