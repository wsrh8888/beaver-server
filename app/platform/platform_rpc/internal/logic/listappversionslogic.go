package logic

import (
	"context"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListAppVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAppVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAppVersionsLogic {
	return &ListAppVersionsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListAppVersionsLogic) ListAppVersions(in *platform_rpc.ListAppVersionsReq) (*platform_rpc.ListAppVersionsRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	db := l.svcCtx.DB.Model(&platform_models.UpdateArchitecture{}).Where("is_active = ?", true)
	if in.AppId != "" {
		db = db.Where("app_id = ?", in.AppId)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计架构失败: %v", err)
		return nil, status.Error(codes.Internal, "获取架构总数失败")
	}

	var architectures []platform_models.UpdateArchitecture
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&architectures).Error; err != nil {
		l.Errorf("查询架构失败: %v", err)
		return nil, status.Error(codes.Internal, "获取架构列表失败")
	}

	items := make([]*platform_rpc.AppVersionsArchItem, 0, len(architectures))
	for _, arch := range architectures {
		var versions []platform_models.UpdateVersion
		if err := l.svcCtx.DB.Where("architecture_id = ?", arch.Id).Order("created_at DESC").Find(&versions).Error; err != nil {
			l.Errorf("查询架构版本失败 arch=%d: %v", arch.Id, err)
			continue
		}

		briefs := make([]*platform_rpc.AppVersionBrief, 0, len(versions))
		for _, ver := range versions {
			briefs = append(briefs, &platform_rpc.AppVersionBrief{
				VersionId: uint64(ver.Id),
				Version:   ver.Version,
			})
		}

		items = append(items, &platform_rpc.AppVersionsArchItem{
			ArchitectureId: uint64(arch.Id),
			ArchId:         uint32(arch.ArchID),
			Description:    arch.Description,
			Versions:       briefs,
		})
	}

	return &platform_rpc.ListAppVersionsRes{Total: total, Architectures: items}, nil
}
