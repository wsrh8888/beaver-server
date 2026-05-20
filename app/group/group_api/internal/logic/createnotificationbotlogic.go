package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateNotificationBotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 在群里创建通知机器人（群管理员操作，返回 Webhook URL + Secret）
func NewCreateNotificationBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNotificationBotLogic {
	return &CreateNotificationBotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNotificationBotLogic) CreateNotificationBot(req *types.CreateNotificationBotReq) (resp *types.CreateNotificationBotRes, err error) {
	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", req.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可创建通知机器人")
	}

	tokenBytes := make([]byte, 32)
	if _, err = rand.Read(tokenBytes); err != nil {
		return nil, errors.New("生成 token 失败")
	}
	secretBytes := make([]byte, 32)
	if _, err = rand.Read(secretBytes); err != nil {
		return nil, errors.New("生成 secret 失败")
	}
	token := hex.EncodeToString(tokenBytes)
	secret := hex.EncodeToString(secretBytes)

	record := open_models.OpenIncomingWebhook{
		Token:     token,
		Secret:    secret,
		AppID:     "GROUP_NOTIFICATION",
		GroupID:   req.GroupID,
		BotUserID: "NOTIFICATION_BOT",
		Name:      req.Name,
		Status:    1,
	}
	if err = l.svcCtx.DB.Create(&record).Error; err != nil {
		return nil, errors.New("创建失败")
	}

	webhookURL := fmt.Sprintf("%s/api/open/v1/webhook/incoming?token=%s", l.svcCtx.Config.ApiBaseUrl, token)
	return &types.CreateNotificationBotRes{
		ID:         int64(record.ID),
		WebhookURL: webhookURL,
		Secret:     secret,
	}, nil
}
