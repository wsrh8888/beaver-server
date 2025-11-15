package logic

import (
	"context"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/feedback/feedback_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeedbackListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取反馈列表
func NewGetFeedbackListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedbackListLogic {
	return &GetFeedbackListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFeedbackListLogic) GetFeedbackList(req *types.GetFeedbackListReq) (resp *types.GetFeedbackListRes, err error) {
	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 构建查询条件
	query := l.svcCtx.DB.Model(&feedback_models.FeedbackModel{})

	// 状态筛选
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 类型筛选
	if req.Type > 0 {
		query = query.Where("type = ?", req.Type)
	}

	// 用户ID筛选
	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
	}

	// 关键词搜索
	if req.Keywords != "" {
		query = query.Where("content LIKE ?", "%"+req.Keywords+"%")
	}

	// 查询总数
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		logx.Errorf("查询反馈总数失败: %v", err)
		return nil, err
	}

	// 查询列表
	var feedbacks []feedback_models.FeedbackModel
	err = query.Order("created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&feedbacks).Error
	if err != nil {
		logx.Errorf("查询反馈列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	list := make([]types.GetFeedbackListItem, len(feedbacks))
	for i, feedback := range feedbacks {
		var handleTime string
		if feedback.HandleTime != nil {
			handleTime = feedback.HandleTime.Format(time.RFC3339)
		}

		list[i] = types.GetFeedbackListItem{
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
		}
	}

	return &types.GetFeedbackListRes{
		List:  list,
		Total: total,
	}, nil
}
