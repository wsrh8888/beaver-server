package logic

import (
	"context"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListReleasePoliciesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListReleasePoliciesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListReleasePoliciesLogic {
	return &ListReleasePoliciesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListReleasePoliciesLogic) ListReleasePolicies(in *platform_rpc.ListReleasePoliciesReq) (*platform_rpc.ListReleasePoliciesRes, error) {
	var list []platform_models.UpdateReleasePolicy
	db := l.svcCtx.DB.Model(&platform_models.UpdateReleasePolicy{})
	if in.AppId != "" {
		db = db.Where("app_id = ?", in.AppId)
	}
	if err := db.Order("architecture_id ASC").Find(&list).Error; err != nil {
		return nil, err
	}

	items := make([]*platform_rpc.ReleasePolicyItem, 0, len(list))
	for _, p := range list {
		items = append(items, &platform_rpc.ReleasePolicyItem{
			Id:              uint64(p.Id),
			AppId:           p.AppID,
			ArchitectureId:    uint64(p.ArchitectureID),
			StableVersionId: uint64(p.StableVersionID),
			GrayVersionId:   uint64(p.GrayVersionID),
			RolloutPercent:    uint32(p.RolloutPercent),
			MinVersion:        p.MinVersion,
			ForceUpdate:       p.ForceUpdate,
			IsActive:          p.IsActive,
			StableVersion:     versionLabel(l.svcCtx.DB, p.StableVersionID),
			GrayVersion:       versionLabel(l.svcCtx.DB, p.GrayVersionID),
			CreatedAt:         p.CreatedAt.String(),
			UpdatedAt:         p.UpdatedAt.String(),
		})
	}
	return &platform_rpc.ListReleasePoliciesRes{Policies: items}, nil
}
