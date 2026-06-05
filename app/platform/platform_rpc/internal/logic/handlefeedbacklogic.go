package logic

import (
	"context"
	"errors"
	"strconv"
	"time"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type HandleFeedbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleFeedbackLogic {
	return &HandleFeedbackLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *HandleFeedbackLogic) HandleFeedback(in *platform_rpc.HandleFeedbackReq) (*platform_rpc.HandleFeedbackRes, error) {
	if in.Status < 1 || in.Status > 4 {
		return nil, status.Error(codes.InvalidArgument, "无效的状态值")
	}

	var feedback platform_models.FeedbackModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&feedback).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "反馈记录不存在")
		}
		l.Errorf("查询反馈失败: %v", err)
		return nil, err
	}

	handlerID, _ := strconv.ParseInt(in.HandlerId, 10, 64)
	now := time.Now()
	if err := l.svcCtx.DB.Model(&feedback).Updates(map[string]interface{}{
		"status":        platform_models.FeedbackStatus(in.Status),
		"handle_result": in.HandleResult,
		"handler_id":    handlerID,
		"handle_time":   &now,
		"updated_at":    now,
	}).Error; err != nil {
		l.Errorf("处理反馈失败: %v", err)
		return nil, status.Error(codes.Internal, "处理反馈失败")
	}

	return &platform_rpc.HandleFeedbackRes{}, nil
}
