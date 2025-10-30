package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/feedback/feedback_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type HandleFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 处理反馈
func NewHandleFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleFeedbackLogic {
	return &HandleFeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandleFeedbackLogic) HandleFeedback(req *types.HandleFeedbackReq) (resp *types.HandleFeedbackRes, err error) {
	// 验证状态值
	if req.Status < 1 || req.Status > 4 {
		return nil, errors.New("无效的状态值")
	}

	// 查询反馈记录
	var feedback feedback_models.FeedbackModel
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&feedback).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("反馈记录不存在: %d", req.Id)
			return nil, errors.New("反馈记录不存在")
		}
		logx.Errorf("查询反馈记录失败: %v", err)
		return nil, err
	}

	// 更新反馈状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":        feedback_models.FeedbackStatus(req.Status),
		"handle_result": req.HandleResult,
		"handler_id":    req.UserID,
		"handle_time":   &now,
		"updated_at":    now,
	}

	err = l.svcCtx.DB.Model(&feedback).Updates(updates).Error
	if err != nil {
		logx.Errorf("更新反馈状态失败: %v", err)
		return nil, errors.New("处理反馈失败")
	}

	return &types.HandleFeedbackRes{}, nil
}
