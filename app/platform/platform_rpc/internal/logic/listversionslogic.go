package logic

import (
	"context"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListVersionsLogic {
	return &ListVersionsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListVersionsLogic) ListVersions(in *platform_rpc.ListVersionsReq) (*platform_rpc.ListVersionsRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	db := l.svcCtx.DB.Model(&platform_models.UpdateVersion{})
	if in.ArchitectureId > 0 {
		db = db.Where("architecture_id = ?", in.ArchitectureId)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计版本失败: %v", err)
		return nil, err
	}

	var list []platform_models.UpdateVersion
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error; err != nil {
		l.Errorf("查询版本列表失败: %v", err)
		return nil, err
	}

	items := make([]*platform_rpc.VersionItem, 0, len(list))
	for _, ver := range list {
		items = append(items, &platform_rpc.VersionItem{
			VersionId:      uint64(ver.Id),
			ArchitectureId: uint64(ver.ArchitectureID),
			Version:        ver.Version,
			FileKey:        ver.FileKey,
			Description:    ver.Description,
			ReleaseNotes:   ver.ReleaseNotes,
			CreatedAt:      ver.CreatedAt.String(),
			UpdatedAt:      ver.UpdatedAt.String(),
		})
	}

	return &platform_rpc.ListVersionsRes{Total: total, Versions: items}, nil
}
