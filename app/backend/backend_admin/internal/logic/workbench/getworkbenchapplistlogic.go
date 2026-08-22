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
	in := &platform_rpc.ListWorkbenchAppsReq{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		Status:   int32(req.Status),
		Keywords: req.Keywords,
	}
	if req.Category != nil {
		v := int32(*req.Category)
		in.Category = &v
	}

	rpcRes, err := l.svcCtx.PlatformRpc.ListWorkbenchApps(l.ctx, in)
	if err != nil {
		l.Errorf("获取工作台应用列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetWorkbenchAppListItem, 0, len(rpcRes.List))
	for _, item := range rpcRes.List {
		list = append(list, toAdminAppItem(item))
	}

	return &types.GetWorkbenchAppListRes{Total: rpcRes.Total, List: list}, nil
}
