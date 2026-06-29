package workbench

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkbenchAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkbenchAppLogic {
	return &GetWorkbenchAppLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetWorkbenchAppLogic) GetWorkbenchApp(req *types.GetWorkbenchAppReq) (*types.GetWorkbenchAppRes, error) {
	rpcRes, err := l.svcCtx.PlatformRpc.GetWorkbenchApp(l.ctx, &platform_rpc.GetWorkbenchAppReq{
		WorkbenchAppId: req.WorkbenchAppID,
	})
	if err != nil {
		l.Errorf("获取工作台应用详情失败: %v", err)
		return nil, err
	}

	app := rpcRes.App
	return &types.GetWorkbenchAppRes{
		WorkbenchAppID: app.WorkbenchAppId,
		Name:           app.Name,
		Description:    app.Description,
		Icon:           app.Icon,
		EntryURL:       app.EntryUrl,
		Category:       app.Category,
		Sort:           int(app.Sort),
		Status:         int(app.Status),
		Remark:         app.Remark,
		CreatedBy:      app.CreatedBy,
		LastModifiedBy: app.LastModifiedBy,
		CreatedAt:      app.CreatedAt,
		UpdatedAt:      app.UpdatedAt,
	}, nil
}
