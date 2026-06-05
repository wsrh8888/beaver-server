package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVersionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetVersionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVersionListLogic {
	return &GetVersionListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetVersionListLogic) GetVersionList(req *types.GetVersionListReq) (resp *types.GetVersionListRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.ListVersions(l.ctx, &platform_rpc.ListVersionsReq{
		ArchitectureId: uint64(req.ArchitectureID),
		Page:           int32(req.Page),
		PageSize:       int32(req.PageSize),
	})
	if err != nil {
		l.Errorf("获取版本列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetVersionListItem, 0, len(rpcRes.Versions))
	for _, ver := range rpcRes.Versions {
		list = append(list, types.GetVersionListItem{
			VersionID:      uint(ver.VersionId),
			ArchitectureID: uint(ver.ArchitectureId),
			Version:        ver.Version,
			FileKey:        ver.FileKey,
			Description:    ver.Description,
			ReleaseNotes:   ver.ReleaseNotes,
			CreatedAt:      ver.CreatedAt,
			UpdatedAt:      ver.UpdatedAt,
		})
	}

	return &types.GetVersionListRes{Total: rpcRes.Total, Versions: list}, nil
}
