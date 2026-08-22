package workbench

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWorkbenchAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkbenchAppLogic {
	return &DeleteWorkbenchAppLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteWorkbenchAppLogic) DeleteWorkbenchApp(req *types.DeleteWorkbenchAppReq) (*types.DeleteWorkbenchAppRes, error) {
	_, err := l.svcCtx.PlatformRpc.DeleteWorkbenchApp(l.ctx, &platform_rpc.DeleteWorkbenchAppReq{
		WorkbenchAppId: req.WorkbenchAppID,
	})
	if err != nil {
		l.Errorf("删除工作台应用失败: %v", err)
		return nil, err
	}

	return &types.DeleteWorkbenchAppRes{}, nil
}
