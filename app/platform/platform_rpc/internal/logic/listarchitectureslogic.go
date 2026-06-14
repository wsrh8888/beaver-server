package logic

import (
	"context"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArchitecturesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListArchitecturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArchitecturesLogic {
	return &ListArchitecturesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListArchitecturesLogic) ListArchitectures(in *platform_rpc.ListArchitecturesReq) (*platform_rpc.ListArchitecturesRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	db := l.svcCtx.DB.Model(&platform_models.UpdateArchitecture{})
	if in.AppId != "" {
		db = db.Where("app_id = ?", in.AppId)
	}
	if in.IsActive != nil && *in.IsActive {
		db = db.Where("is_active = ?", true)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计架构失败: %v", err)
		return nil, err
	}

	var list []platform_models.UpdateArchitecture
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error; err != nil {
		l.Errorf("查询架构列表失败: %v", err)
		return nil, err
	}

	appNameMap := l.loadAppNameMap(list)

	items := make([]*platform_rpc.ArchitectureItem, 0, len(list))
	for _, arch := range list {
		items = append(items, &platform_rpc.ArchitectureItem{
			Id:          uint64(arch.Id),
			AppId:       arch.AppID,
			AppName:     appNameMap[arch.AppID],
			PlatformId:  uint32(arch.PlatformID),
			ArchId:      uint32(arch.ArchID),
			Description: arch.Description,
			IsActive:    arch.IsActive,
			CreatedAt:   arch.CreatedAt.String(),
			UpdatedAt:   arch.UpdatedAt.String(),
		})
	}

	return &platform_rpc.ListArchitecturesRes{Total: total, Architectures: items}, nil
}

func (l *ListArchitecturesLogic) loadAppNameMap(list []platform_models.UpdateArchitecture) map[string]string {
	appNameMap := make(map[string]string)
	if len(list) == 0 {
		return appNameMap
	}

	appIDs := make([]string, 0, len(list))
	seen := make(map[string]struct{}, len(list))
	for _, arch := range list {
		if arch.AppID == "" {
			continue
		}
		if _, ok := seen[arch.AppID]; ok {
			continue
		}
		seen[arch.AppID] = struct{}{}
		appIDs = append(appIDs, arch.AppID)
	}
	if len(appIDs) == 0 {
		return appNameMap
	}

	var apps []platform_models.UpdateApp
	if err := l.svcCtx.DB.Where("app_id IN ?", appIDs).Find(&apps).Error; err != nil {
		l.Errorf("查询应用名称失败: %v", err)
		return appNameMap
	}
	for _, app := range apps {
		appNameMap[app.AppID] = app.Name
	}
	return appNameMap
}
