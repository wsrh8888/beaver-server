package event

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEventLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventLogsLogic {
	return &GetEventLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventLogsLogic) GetEventLogs(req *types.GetEventLogsReq) (resp *types.GetEventLogsRes, err error) {
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	query := l.svcCtx.DB.Model(&open_models.OpenWebhookLog{}).Where("app_id = ?", req.AppID)
	if req.SubscriptionID > 0 {
		query = query.Where("subscription_id = ?", req.SubscriptionID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询日志失败")
	}

	var logs []open_models.OpenWebhookLog
	if err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, errors.New("查询日志失败")
	}

	list := make([]types.GetEventLogsResItem, 0, len(logs))
	for _, item := range logs {
		list = append(list, types.GetEventLogsResItem{
			ID:           uint64(item.ID),
			EventID:      item.EventID,
			EventType:    item.EventType,
			TargetURL:    item.TargetURL,
			ResponseCode: item.HTTPStatus,
			CostMs:       item.LatencyMs,
			RetryCount:   item.RetryCount,
			Status:       item.Status,
			ErrorMsg:     item.ErrorMessage,
			CreatedAt:    item.CreatedAt.Unix(),
		})
	}

	return &types.GetEventLogsRes{
		Total: total,
		List:  list,
	}, nil
}
