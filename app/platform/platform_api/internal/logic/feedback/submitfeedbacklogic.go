package feedback

import (
	"context"
	"errors"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/app/platform/platform_models"

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
