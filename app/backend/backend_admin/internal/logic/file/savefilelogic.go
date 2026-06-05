package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveFileLogic {
	return &SaveFileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SaveFileLogic) SaveFile(req *types.SaveFileReq) (resp *types.SaveFileRes, err error) {
	rpcRes, err := l.svcCtx.FileRpc.SaveFile(l.ctx, &file_rpc.SaveFileReq{
		OriginalName: req.OriginalName,
		Size:         req.Size,
		Path:         req.Path,
		Md5:          req.Md5,
		Type:         req.Type,
		Source:       req.Source,
		FileInfoJson: req.FileInfo,
	})
	if err != nil {
		l.Errorf("保存文件失败: %v", err)
		return nil, err
	}
	return &types.SaveFileRes{FileKey: rpcRes.FileKey}, nil
}
