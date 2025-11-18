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

type GetFeedbackDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取反馈详情
func NewGetFeedbackDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedbackDetailLogic {
	return &GetFeedbackDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFeedbackDetailLogic) GetFeedbackDetail(req *types.GetFeedbackDetailReq) (resp *types.GetFeedbackDetailRes, err error) {
	var feedback feedback_models.FeedbackModel
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&feedback).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("反馈记录不存在: %d", req.Id)
			return nil, errors.New("反馈记录不存在")
		}
		logx.Errorf("查询反馈详情失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var handleTime string
	if feedback.HandleTime != nil {
		handleTime = feedback.HandleTime.Format(time.RFC3339)
	}

	return &types.GetFeedbackDetailRes{
		Id:           feedback.Id,
		UserId:       feedback.UserID,
		Content:      feedback.Content,
		Type:         int(feedback.Type),
		Status:       int(feedback.Status),
		FileNames:    []string(feedback.FileNames),
		HandlerId:    feedback.HandlerID,
		HandleTime:   handleTime,
		HandleResult: feedback.HandleResult,
		CreatedAt:    time.Time(feedback.CreatedAt).Format(time.RFC3339),
		UpdatedAt:    time.Time(feedback.UpdatedAt).Format(time.RFC3339),
	}, nil
}
