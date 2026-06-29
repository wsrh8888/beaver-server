package workbench

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWorkbenchAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkbenchAppLogic {
	return &CreateWorkbenchAppLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateWorkbenchAppLogic) CreateWorkbenchApp(req *types.CreateWorkbenchAppReq) (*types.CreateWorkbenchAppRes, error) {
	rpcRes, err := l.svcCtx.PlatformRpc.CreateWorkbenchApp(l.ctx, &platform_rpc.CreateWorkbenchAppReq{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		EntryUrl:    req.EntryURL,
		Category:    req.Category,
		Sort:        int32(req.Sort),
		Status:      int32(req.Status),
		Remark:      req.Remark,
		OperatorId:  req.UserID,
	})
	if err != nil {
		l.Errorf("创建工作台应用失败: %v", err)
		return nil, err
	}

	return &types.CreateWorkbenchAppRes{WorkbenchAppID: rpcRes.WorkbenchAppId}, nil
}
