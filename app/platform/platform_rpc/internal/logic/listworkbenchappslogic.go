package logic

import (
	"context"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkbenchAppsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListWorkbenchAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkbenchAppsLogic {
	return &ListWorkbenchAppsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListWorkbenchAppsLogic) ListWorkbenchApps(in *platform_rpc.ListWorkbenchAppsReq) (*platform_rpc.ListWorkbenchAppsRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := l.svcCtx.DB.Model(&platform_models.WorkbenchApp{})
	if in.Status > 0 {
		db = db.Where("status = ?", in.Status)
	}
	if in.Category != "" {
		db = db.Where("category = ?", in.Category)
	}
	if in.Keywords != "" {
		like := "%" + in.Keywords + "%"
		db = db.Where("name LIKE ? OR description LIKE ?", like, like)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计工作台应用失败: %v", err)
		return nil, err
	}

	var list []platform_models.WorkbenchApp
	if err := db.Order("sort ASC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.Errorf("查询工作台应用列表失败: %v", err)
		return nil, err
	}

	items := make([]*platform_rpc.WorkbenchAppItem, 0, len(list))
	for _, app := range list {
		items = append(items, toWorkbenchAppItem(app))
	}

	return &platform_rpc.ListWorkbenchAppsRes{Total: total, List: items}, nil
}
