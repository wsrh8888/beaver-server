package workbench

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWorkbenchAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkbenchAppLogic {
	return &UpdateWorkbenchAppLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateWorkbenchAppLogic) UpdateWorkbenchApp(req *types.UpdateWorkbenchAppReq) (*types.UpdateWorkbenchAppRes, error) {
	in := &platform_rpc.UpdateWorkbenchAppReq{
		WorkbenchAppId: req.WorkbenchAppID,
		Name:           req.Name,
		Description:    req.Description,
		Icon:           req.Icon,
		EntryUrl:       req.EntryURL,
		Category:       req.Category,
		Remark:         req.Remark,
		OperatorId:     req.UserID,
	}
	if req.Sort != nil {
		sortVal := int32(*req.Sort)
		in.Sort = &sortVal
	}
	if req.Status != nil {
		statusVal := int32(*req.Status)
		in.Status = &statusVal
	}

	_, err := l.svcCtx.PlatformRpc.UpdateWorkbenchApp(l.ctx, in)
	if err != nil {
		l.Errorf("更新工作台应用失败: %v", err)
		return nil, err
	}

	return &types.UpdateWorkbenchAppRes{}, nil
}
