package logic

import (
	"context"

	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadQiniuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadQiniuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadQiniuLogic {
	return &FileUploadQiniuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadQiniuLogic) FileUploadQiniu(req *types.FileReq) (resp *types.FileRes, err error) {
	// todo: add your logic here and delete this line

	return &types.FileRes{}, nil
}
