package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteFileLogic {
	return &BatchDeleteFileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *BatchDeleteFileLogic) BatchDeleteFile(req *types.BatchDeleteFileReq) (resp *types.BatchDeleteFileRes, err error) {
	ids := make([]uint64, 0, len(req.Ids))
	for _, id := range req.Ids {
		ids = append(ids, uint64(id))
	}

	_, err = l.svcCtx.FileRpc.BatchDeleteFiles(l.ctx, &file_rpc.BatchDeleteFilesReq{Ids: ids})
	if err != nil {
		l.Errorf("批量删除文件失败: %v", err)
		return nil, err
	}
	return &types.BatchDeleteFileRes{}, nil
}
