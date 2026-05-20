package bot

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	models "beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateIncomingWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建 Incoming Webhook（开放平台应用维度）
func NewCreateIncomingWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateIncomingWebhookLogic {
	return &CreateIncomingWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateIncomingWebhookLogic) CreateIncomingWebhook(req *types.CreateIncomingWebhookReq) (resp *types.CreateIncomingWebhookRes, err error) {
	tokenBytes := make([]byte, 32)
	if _, err = rand.Read(tokenBytes); err != nil {
		return nil, errors.New("生成 access_token 失败")
	}
	secretBytes := make([]byte, 32)
	if _, err = rand.Read(secretBytes); err != nil {
		return nil, errors.New("生成 secret 失败")
	}
	token := hex.EncodeToString(tokenBytes)
	secret := hex.EncodeToString(secretBytes)

	botUserID := "bot_" + req.AppID
	webhook := models.OpenIncomingWebhook{
		Token:     token,
		Secret:    secret,
		AppID:     req.AppID,
		GroupID:   req.GroupID,
		BotUserID: botUserID,
		Name:      req.Name,
		Status:    1,
	}

	if err := l.svcCtx.DB.Create(&webhook).Error; err != nil {
		return nil, errors.New("创建 Incoming Webhook 失败")
	}

	webhookURL := fmt.Sprintf("%s/api/open/v1/webhook/incoming?access_token=%s", l.svcCtx.Config.ApiBaseUrl, token)

	return &types.CreateIncomingWebhookRes{
		Webhook: types.IncomingWebhookInfo{
			ID:         fmt.Sprintf("%d", webhook.ID),
			Token:      webhook.Token,
			Secret:     secret,
			AppID:      webhook.AppID,
			GroupID:    webhook.GroupID,
			BotUserID:  webhook.BotUserID,
			Name:       webhook.Name,
			WebhookURL: webhookURL,
			Status:     webhook.Status,
			CreatedAt:  time.Now().Unix(),
		},
	}, nil
}
