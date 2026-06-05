package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHandleFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleFeedbackLogic {
	return &HandleFeedbackLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *HandleFeedbackLogic) HandleFeedback(req *types.HandleFeedbackReq) (resp *types.HandleFeedbackRes, err error) {
	_, err = l.svcCtx.PlatformRpc.HandleFeedback(l.ctx, &platform_rpc.HandleFeedbackReq{
		Id:           uint64(req.Id),
		Status:       int32(req.Status),
		HandleResult: req.HandleResult,
		HandlerId:    req.UserID,
	})
	if err != nil {
		l.Errorf("处理反馈失败: %v", err)
		return nil, err
	}
	return &types.HandleFeedbackRes{}, nil
}
