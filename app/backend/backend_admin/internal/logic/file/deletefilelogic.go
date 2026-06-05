package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteFileLogic) DeleteFile(req *types.DeleteFileReq) (resp *types.DeleteFileRes, err error) {
	_, err = l.svcCtx.FileRpc.DeleteFile(l.ctx, &file_rpc.DeleteFileReq{Id: uint64(req.Id)})
	if err != nil {
		l.Errorf("删除文件失败: %v", err)
		return nil, err
	}
	return &types.DeleteFileRes{}, nil
}
