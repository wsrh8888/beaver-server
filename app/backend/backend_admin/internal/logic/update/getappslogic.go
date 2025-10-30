package logic

import (
	"beaver/app/update/update_models"
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取应用列表
func NewGetAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppsLogic {
	return &GetAppsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppsLogic) GetApps(req *types.GetAppsReq) (resp *types.GetAppsRes, err error) {
	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 构建查询
	query := l.svcCtx.DB.Model(&update_models.UpdateApp{})
	if req.IsActive {
		query = query.Where("is_active = ?", true)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logx.Errorf("Failed to count apps: %v", err)
		return nil, err
	}

	// 获取分页数据
	var apps []update_models.UpdateApp
	if err := query.Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order("created_at DESC").
		Find(&apps).Error; err != nil {
		logx.Errorf("Failed to get apps: %v", err)
		return nil, err
	}

	// 构建响应
	appList := make([]types.AppInfo, 0, len(apps))
	for _, app := range apps {
		appList = append(appList, types.AppInfo{
			Id:          app.Id,
			AppID:       app.UUID,
			Name:        app.Name,
			Description: app.Description,
			IsActive:    app.IsActive,
			CreatedAt:   app.CreatedAt.String(),
			UpdatedAt:   app.UpdatedAt.String(),
		})
	}

	return &types.GetAppsRes{
		Total: total,
		Apps:  appList,
	}, nil
}
