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

type GetReleasePoliciesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetReleasePoliciesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReleasePoliciesLogic {
	return &GetReleasePoliciesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReleasePoliciesLogic) GetReleasePolicies(req *types.GetReleasePoliciesReq) (resp *types.GetReleasePoliciesRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.ListReleasePolicies(l.ctx, &platform_rpc.ListReleasePoliciesReq{
		AppId: req.AppID,
	})
	if err != nil {
		l.Errorf("获取发版策略失败: %v", err)
		return nil, err
	}

	policies := make([]types.ReleasePolicyItem, 0, len(rpcRes.Policies))
	for _, p := range rpcRes.Policies {
		policies = append(policies, types.ReleasePolicyItem{
			ID:              uint(p.Id),
			AppID:           p.AppId,
			ArchitectureID:  uint(p.ArchitectureId),
			StableVersionID: uint(p.StableVersionId),
			GrayVersionID:   uint(p.GrayVersionId),
			RolloutPercent:  uint(p.RolloutPercent),
			MinVersion:      p.MinVersion,
			ForceUpdate:     p.ForceUpdate,
			IsActive:        p.IsActive,
			StableVersion:   p.StableVersion,
			GrayVersion:     p.GrayVersion,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
		})
	}
	return &types.GetReleasePoliciesRes{Policies: policies}, nil
}
