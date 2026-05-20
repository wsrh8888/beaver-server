package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
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
	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", req.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可查看通知机器人")
	}

	var records []open_models.OpenIncomingWebhook
	if err = l.svcCtx.DB.Where("group_id = ? AND app_id = ?", req.GroupID, "GROUP_NOTIFICATION").
		Order("id DESC").Find(&records).Error; err != nil {
		return nil, err
	}

	baseURL := l.svcCtx.Config.ApiBaseUrl
	items := make([]types.ListNotificationBotsItem, 0, len(records))
	for _, w := range records {
		items = append(items, types.ListNotificationBotsItem{
			ID:            int64(w.ID),
			Name:          w.Name,
			Description:   w.Description,
			Avatar:        w.Avatar,
			WebhookURL:    fmt.Sprintf("%s/api/open/v1/webhook/incoming?access_token=%s", baseURL, w.Token),
			Status:        w.Status,
			CreatorUserID: w.CreatorUserID,
			CreatedAt:     w.CreatedAt.Unix(),
		})
	}

	return &types.ListNotificationBotsRes{List: items}, nil
}
