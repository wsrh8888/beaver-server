package logic

import (
	"context"
	"errors"

	"beaver/app/feedback/feedback_api/internal/svc"
	"beaver/app/feedback/feedback_api/internal/types"
	"beaver/app/feedback/feedback_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 提交反馈
func NewSubmitFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitFeedbackLogic {
	return &SubmitFeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitFeedbackLogic) SubmitFeedback(req *types.SubmitFeedbackReq) (resp *types.SubmitFeedbackRes, err error) {
	// 创建反馈记录
	feedback := &feedback_models.FeedbackModel{
		UserID:    req.UserID,
		Content:   req.Content,
		Type:      feedback_models.FeedbackType(req.Type),
		Status:    feedback_models.FeedbackStatusPending,
		FileNames: feedback_models.FileNames(req.FileNames),
	}

	// 保存到数据库
	err = l.svcCtx.DB.Create(feedback).Error
	if err != nil {
		logx.Errorf("创建反馈失败: %v", err)
		return nil, errors.New("提交反馈失败")
	}

	return &types.SubmitFeedbackRes{}, nil
}
