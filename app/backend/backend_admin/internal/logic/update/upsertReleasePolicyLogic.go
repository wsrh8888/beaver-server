// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpsertReleasePolicyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpsertReleasePolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpsertReleasePolicyLogic {
	return &UpsertReleasePolicyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpsertReleasePolicyLogic) UpsertReleasePolicy(req *types.UpsertReleasePolicyReq) (resp *types.UpsertReleasePolicyRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.UpsertReleasePolicy(l.ctx, &platform_rpc.UpsertReleasePolicyReq{
		AppId:           req.AppID,
		ArchitectureId:  uint64(req.ArchitectureID),
		StableVersionId: uint64(req.StableVersionID),
		GrayVersionId:   uint64(req.GrayVersionID),
		RolloutPercent:  uint32(req.RolloutPercent),
		MinVersion:      req.MinVersion,
		ForceUpdate:     req.ForceUpdate,
		IsActive:        req.IsActive,
	})
	if err != nil {
		l.Errorf("保存发版策略失败: %v", err)
		return nil, err
	}
	return &types.UpsertReleasePolicyRes{ID: uint(rpcRes.Id)}, nil
}
