package feedback

import (
	"context"
	"errors"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"

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
	feedback := &platform_models.FeedbackModel{
		UserID:    req.UserID,
		Content:   req.Content,
		Type:      platform_models.FeedbackType(req.Type),
		Status:    platform_models.FeedbackStatusPending,
		FileNames: platform_models.FileNames(req.FileNames),
	}

	if err := l.svcCtx.DB.Create(feedback).Error; err != nil {
		logx.Errorf("create feedback failed: %v", err)
		return nil, errors.New("提交反馈失败")
	}

	return &types.SubmitFeedbackRes{}, nil
}
