package logic

import (
	"context"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubmitFeedbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitFeedbackLogic {
	return &SubmitFeedbackLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SubmitFeedbackLogic) SubmitFeedback(in *platform_rpc.SubmitFeedbackReq) (*platform_rpc.SubmitFeedbackRes, error) {
	if in.UserId == "" || in.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "用户ID和反馈内容不能为空")
	}

	feedback := platform_models.FeedbackModel{
		UserID:    in.UserId,
		Content:   in.Content,
		Type:      platform_models.FeedbackType(in.Type),
		Status:    platform_models.FeedbackStatusPending,
		FileNames: platform_models.FileNames(in.FileNames),
	}
	if err := l.svcCtx.DB.Create(&feedback).Error; err != nil {
		l.Errorf("提交反馈失败: %v", err)
		return nil, status.Error(codes.Internal, "提交反馈失败")
	}

	return &platform_rpc.SubmitFeedbackRes{Id: uint64(feedback.Id)}, nil
}
