package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/feedback/feedback_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除反馈
func NewDeleteFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFeedbackLogic {
	return &DeleteFeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFeedbackLogic) DeleteFeedback(req *types.DeleteFeedbackReq) (resp *types.DeleteFeedbackRes, err error) {
	// 检查反馈是否存在
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

	// 软删除反馈记录
	err = l.svcCtx.DB.Delete(&feedback).Error
	if err != nil {
		logx.Errorf("删除反馈记录失败: %v", err)
		return nil, errors.New("删除反馈失败")
	}

	return &types.DeleteFeedbackRes{}, nil
}
