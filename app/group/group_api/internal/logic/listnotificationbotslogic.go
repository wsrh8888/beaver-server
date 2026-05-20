package logic

import (
	"context"
	"fmt"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNotificationBotsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取群内所有通知机器人列表
func NewListNotificationBotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotificationBotsLogic {
	return &ListNotificationBotsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListNotificationBotsLogic) ListNotificationBots(req *types.ListNotificationBotsReq) (resp *types.ListNotificationBotsRes, err error) {
	var records []open_models.OpenIncomingWebhook
	if err = l.svcCtx.DB.Where("group_id = ? AND app_id = ?", req.GroupID, "GROUP_NOTIFICATION").
		Order("id DESC").Find(&records).Error; err != nil {
		return nil, err
	}

	baseURL := l.svcCtx.Config.ApiBaseUrl
	items := make([]types.NotificationBotItem, 0, len(records))
	for _, w := range records {
		items = append(items, types.NotificationBotItem{
			ID:         int64(w.ID),
			Name:       w.Name,
			WebhookURL: fmt.Sprintf("%s/api/open/v1/webhook/incoming?token=%s", baseURL, w.Token),
			CreatedAt:  w.CreatedAt.Unix(),
		})
	}

	return &types.ListNotificationBotsRes{List: items}, nil
}
