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

type DeleteFeedbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFeedbackLogic {
	return &DeleteFeedbackLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteFeedbackLogic) DeleteFeedback(in *platform_rpc.DeleteFeedbackReq) (*platform_rpc.DeleteFeedbackRes, error) {
	var feedback platform_models.FeedbackModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&feedback).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "反馈记录不存在")
		}
		l.Errorf("查询反馈失败: %v", err)
		return nil, err
	}

	if err := l.svcCtx.DB.Delete(&feedback).Error; err != nil {
		l.Errorf("删除反馈失败: %v", err)
		return nil, status.Error(codes.Internal, "删除反馈失败")
	}

	return &platform_rpc.DeleteFeedbackRes{}, nil
}
