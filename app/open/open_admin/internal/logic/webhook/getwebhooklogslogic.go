package webhook

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"
	models "beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWebhookLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 Webhook 日志
func NewGetWebhookLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWebhookLogsLogic {
	return &GetWebhookLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWebhookLogsLogic) GetWebhookLogs(req *types.GetWebhookLogsReq) (resp *types.GetWebhookLogsRes, err error) {
	// 1. 构建查询
	query := l.svcCtx.DB.Model(&models.OpenWebhookLog{}).Where("app_id = ?", req.AppID)
	
	if req.EventType != "" {
		query = query.Where("event_type = ?", req.EventType)
	}

	// 2. 获取总数
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		return nil, errors.New("查询失败")
	}

	// 3. 分页查询
	var logs []models.OpenWebhookLog
	offset := (req.Page - 1) * req.PageSize
	err = query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&logs).Error
	if err != nil {
		return nil, errors.New("查询失败")
	}

	// 4. 转换为响应格式
	list := make([]types.WebhookLogItem, 0, len(logs))
	for _, log := range logs {
		list = append(list, types.WebhookLogItem{
			ID:           fmt.Sprintf("%d", log.ID),
			EventType:    log.EventType,
			Payload:      log.Payload,
			ResponseCode: log.ResponseCode,
			RetryCount:   log.RetryCount,
			Status:       log.Status,
			CreatedAt:    log.CreatedAt.Unix(),
		})
	}

	return &types.GetWebhookLogsRes{
		Total: total,
		List:  list,
	}, nil
}
