package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadQiniuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件上传七牛云
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
