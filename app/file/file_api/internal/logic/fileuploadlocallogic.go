package logic

import (
	"context"

	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadLocalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件上传本地
func NewFileUploadLocalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLocalLogic {
	return &FileUploadLocalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLocalLogic) FileUploadLocal(req *types.FileReq) (resp *types.FileRes, err error) {
	// todo: add your logic here and delete this line

	return &types.FileRes{}, nil
}
