package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddVersionLogic {
	return &AddVersionLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AddVersionLogic) AddVersion(req *types.AddVersionReq) (resp *types.AddVersionRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.CreateVersion(l.ctx, &platform_rpc.CreateVersionReq{
		ArchitectureId: uint64(req.ArchitectureID),
		Version:        req.Version,
		FileUrl:        req.FileUrl,
		Description:    req.Description,
		ReleaseNotes:   req.ReleaseNotes,
	})
	if err != nil {
		l.Errorf("创建版本失败: %v", err)
		return nil, err
	}
	return &types.AddVersionRes{VersionID: uint(rpcRes.VersionId)}, nil
}
