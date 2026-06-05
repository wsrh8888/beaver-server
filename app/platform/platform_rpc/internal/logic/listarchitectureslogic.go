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

	db := l.svcCtx.DB.Model(&platform_models.UpdateArchitecture{}).Preload("App")
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

	items := make([]*platform_rpc.ArchitectureItem, 0, len(list))
	for _, arch := range list {
		appName := ""
		if arch.App != nil {
			appName = arch.App.Name
		}
		items = append(items, &platform_rpc.ArchitectureItem{
			Id:          uint64(arch.Id),
			AppId:       arch.AppID,
			AppName:     appName,
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
