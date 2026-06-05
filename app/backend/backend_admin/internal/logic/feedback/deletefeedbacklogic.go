package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFeedbackLogic {
	return &DeleteFeedbackLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteFeedbackLogic) DeleteFeedback(req *types.DeleteFeedbackReq) (resp *types.DeleteFeedbackRes, err error) {
	_, err = l.svcCtx.PlatformRpc.DeleteFeedback(l.ctx, &platform_rpc.DeleteFeedbackReq{Id: uint64(req.Id)})
	if err != nil {
		l.Errorf("删除反馈失败: %v", err)
		return nil, err
	}
	return &types.DeleteFeedbackRes{}, nil
}
