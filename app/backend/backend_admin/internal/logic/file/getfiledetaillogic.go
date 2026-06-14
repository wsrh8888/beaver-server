package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFileDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileDetailLogic {
	return &GetFileDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetFileDetailLogic) GetFileDetail(req *types.GetFileDetailReq) (resp *types.GetFileDetailRes, err error) {
	rpcRes, err := l.svcCtx.FileRpc.GetFileById(l.ctx, &file_rpc.GetFileByIdReq{Id: uint64(req.Id)})
	if err != nil {
		l.Errorf("获取文件详情失败: %v", err)
		return nil, err
	}

	f := rpcRes.File
	return &types.GetFileDetailRes{
		Id:           uint(f.Id),
		FileName:     f.FileKey,
		OriginalName: f.OriginalName,
		Size:         f.Size,
		Path:         f.Path,
		Md5:          f.Md5,
		Type:         f.Type,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
	}, nil
}
