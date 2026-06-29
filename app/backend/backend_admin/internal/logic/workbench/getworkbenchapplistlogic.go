package workbench

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkbenchAppListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWorkbenchAppListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkbenchAppListLogic {
	return &GetWorkbenchAppListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetWorkbenchAppListLogic) GetWorkbenchAppList(req *types.GetWorkbenchAppListReq) (*types.GetWorkbenchAppListRes, error) {
	rpcRes, err := l.svcCtx.PlatformRpc.ListWorkbenchApps(l.ctx, &platform_rpc.ListWorkbenchAppsReq{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		Status:   int32(req.Status),
		Category: req.Category,
		Keywords: req.Keywords,
	})
	if err != nil {
		l.Errorf("获取工作台应用列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetWorkbenchAppListItem, 0, len(rpcRes.List))
	for _, item := range rpcRes.List {
		list = append(list, types.GetWorkbenchAppListItem{
			WorkbenchAppID: item.WorkbenchAppId,
			Name:           item.Name,
			Description:    item.Description,
			Icon:           item.Icon,
			EntryURL:       item.EntryUrl,
			Category:       item.Category,
			Sort:           int(item.Sort),
			Status:         int(item.Status),
			Remark:         item.Remark,
			CreatedBy:      item.CreatedBy,
			LastModifiedBy: item.LastModifiedBy,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	return &types.GetWorkbenchAppListRes{Total: rpcRes.Total, List: list}, nil
}
