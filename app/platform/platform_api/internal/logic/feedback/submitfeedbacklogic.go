package feedback

import (
	"context"
	"errors"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitFeedbackLogic {
	return &SubmitFeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitFeedbackLogic) SubmitFeedback(req *types.SubmitFeedbackReq) (*types.SubmitFeedbackRes, error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if req.Content == "" {
		return nil, errors.New("反馈内容不能为空")
	}
	if req.Type < 1 || req.Type > 4 {
		return nil, errors.New("反馈类型不合法")
	}

	_, err := l.svcCtx.PlatformRpc.SubmitFeedback(l.ctx, &platform_rpc.SubmitFeedbackReq{
		UserId:    req.UserID,
		Content:   req.Content,
		Type:      int32(req.Type),
		FileNames: req.FileNames,
	})
	if err != nil {
		logx.Errorf("submit feedback failed: %v", err)
		return nil, errors.New("提交反馈失败")
	}

	return &types.SubmitFeedbackRes{}, nil
}
