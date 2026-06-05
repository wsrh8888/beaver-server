package logic

import (
	"context"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWebhookLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListWebhookLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWebhookLogsLogic {
	return &ListWebhookLogsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListWebhookLogsLogic) ListWebhookLogs(in *open_rpc.ListWebhookLogsReq) (*open_rpc.ListWebhookLogsRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := l.svcCtx.DB.Model(&open_models.OpenWebhookLog{})
	if in.AppId != "" {
		db = db.Where("app_id = ?", in.AppId)
	}
	if in.EventType != "" {
		db = db.Where("event_type = ?", in.EventType)
	}
	switch in.Status {
	case 1:
		db = db.Where("status = ?", 1)
	case 2:
		db = db.Where("status = ?", 0)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计 Webhook 日志失败: %v", err)
		return nil, err
	}

	var logs []open_models.OpenWebhookLog
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error; err != nil {
		l.Errorf("查询 Webhook 日志失败: %v", err)
		return nil, err
	}

	list := make([]*open_rpc.WebhookLogItem, 0, len(logs))
	for _, item := range logs {
		list = append(list, &open_rpc.WebhookLogItem{
			Id:           uint64(item.ID),
			AppId:        item.AppID,
			EventId:      item.EventID,
			EventType:    item.EventType,
			TargetUrl:    item.TargetURL,
			HttpStatus:   int32(item.HTTPStatus),
			LatencyMs:    item.LatencyMs,
			RetryCount:   int32(item.RetryCount),
			Status:       int32(item.Status),
			ErrorMessage: item.ErrorMessage,
			CreatedAt:    item.CreatedAt.Unix(),
		})
	}

	return &open_rpc.ListWebhookLogsRes{Total: total, List: list}, nil
}
