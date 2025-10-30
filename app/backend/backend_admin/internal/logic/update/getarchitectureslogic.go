package logic

import (
	"beaver/app/update/update_models"
	"beaver/common/list_query"
	"beaver/common/models"
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArchitecturesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取架构列表
func NewGetArchitecturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArchitecturesLogic {
	return &GetArchitecturesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArchitecturesLogic) GetArchitectures(req *types.GetArchitecturesReq) (resp *types.GetArchitecturesRes, err error) {
	// 构建查询条件
	query := l.svcCtx.DB.Preload("App") // 预加载 App 信息

	// 添加查询条件
	if req.AppID != "" {
		query = query.Where("app_id = ?", req.AppID)
	}
	if req.IsActive {
		query = query.Where("is_active = ?", true)
	}

	// 使用 list_query 进行查询
	list, total, err := list_query.ListQuery(l.svcCtx.DB.Preload("App"), update_models.UpdateArchitecture{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.PageSize,
		},
		Where: query,
	})

	if err != nil {
		logx.Errorf("Failed to get architectures: %v", err)
		return nil, err
	}

	// 构建响应
	architectureList := make([]types.ArchitectureInfo, 0, len(list))
	for _, arch := range list {
		architectureList = append(architectureList, types.ArchitectureInfo{
			Id:          uint(arch.Id),
			AppID:       arch.AppID,
			AppName:     arch.App.Name, // 添加应用名称
			PlatformID:  arch.PlatformID,
			ArchID:      arch.ArchID,
			Description: arch.Description,
			IsActive:    arch.IsActive,
			CreatedAt:   arch.CreatedAt.String(),
			UpdatedAt:   arch.UpdatedAt.String(),
		})
	}

	return &types.GetArchitecturesRes{
		Total:         total,
		Architectures: architectureList,
	}, nil
}
