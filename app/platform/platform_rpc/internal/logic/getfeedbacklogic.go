package logic

import (
	"context"
	"errors"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GetFeedbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedbackLogic {
	return &GetFeedbackLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetFeedbackLogic) GetFeedback(in *platform_rpc.GetFeedbackReq) (*platform_rpc.GetFeedbackRes, error) {
	var feedback platform_models.FeedbackModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&feedback).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "反馈记录不存在")
		}
		l.Errorf("查询反馈详情失败: %v", err)
		return nil, err
	}

	return &platform_rpc.GetFeedbackRes{Feedback: toFeedbackItem(feedback)}, nil
}
