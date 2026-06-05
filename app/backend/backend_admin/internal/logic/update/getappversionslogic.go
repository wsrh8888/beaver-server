package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppVersionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAppVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppVersionsLogic {
	return &GetAppVersionsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAppVersionsLogic) GetAppVersions(req *types.GetAppVersionsReq) (resp *types.GetAppVersionsRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.ListAppVersions(l.ctx, &platform_rpc.ListAppVersionsReq{
		AppId:    req.AppID,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		l.Errorf("获取应用版本失败: %v", err)
		return nil, err
	}

	list := make([]types.GetAppVersionsItem, 0, len(rpcRes.Architectures))
	for _, arch := range rpcRes.Architectures {
		versions := make([]types.GetAppVersionsVersionItem, 0, len(arch.Versions))
		for _, ver := range arch.Versions {
			versions = append(versions, types.GetAppVersionsVersionItem{
				VersionID: uint(ver.VersionId),
				Version:   ver.Version,
			})
		}
		list = append(list, types.GetAppVersionsItem{
			ArchitectureID: uint(arch.ArchitectureId),
			ArchID:         uint(arch.ArchId),
			Description:    arch.Description,
			Versions:       versions,
		})
	}

	return &types.GetAppVersionsRes{Total: rpcRes.Total, Architectures: list}, nil
}
