package logic

import (
	"context"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListEnabledWorkbenchAppsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListEnabledWorkbenchAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEnabledWorkbenchAppsLogic {
	return &ListEnabledWorkbenchAppsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListEnabledWorkbenchAppsLogic) ListEnabledWorkbenchApps(in *platform_rpc.ListEnabledWorkbenchAppsReq) (*platform_rpc.ListEnabledWorkbenchAppsRes, error) {
	db := l.svcCtx.DB.Model(&platform_models.WorkbenchApp{}).Where("status = ?", 1)
	if in.Category != nil {
		db = db.Where("category = ?", *in.Category)
	}
	// client_scope: 0=全部端可见；1/2 仅对应端；请求 1/2 时同时返回 0
	if in.ClientScope == 1 || in.ClientScope == 2 {
		db = db.Where("client_scope IN ?", []int8{0, int8(in.ClientScope)})
	}

	var list []platform_models.WorkbenchApp
	if err := db.Order("sort ASC, id DESC").Find(&list).Error; err != nil {
		l.Errorf("查询上架工作台应用失败: %v", err)
		return nil, err
	}

	items := make([]*platform_rpc.WorkbenchAppPublicItem, 0, len(list))
	for _, app := range list {
		items = append(items, toWorkbenchAppPublicItem(app))
	}

	return &platform_rpc.ListEnabledWorkbenchAppsRes{List: items}, nil
}
