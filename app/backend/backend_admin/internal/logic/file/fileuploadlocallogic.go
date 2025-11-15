package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

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

func (l *FileUploadLocalLogic) FileUploadLocal(req *types.FileUploadLocalReq) (resp *types.FileUploadLocalRes, err error) {
	// Logic层主要处理业务逻辑，现在业务逻辑在Handler中处理
	// 这里返回空的响应，由Handler填充实际数据
	return &types.FileUploadLocalRes{}, nil
}
