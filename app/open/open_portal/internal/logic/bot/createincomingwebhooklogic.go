package bot

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	util "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateIncomingWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建 Incoming Webhook
func NewCreateIncomingWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateIncomingWebhookLogic {
	return &CreateIncomingWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateIncomingWebhookLogic) CreateIncomingWebhook(req *types.CreateIncomingWebhookReq) (resp *types.CreateIncomingWebhookRes, err error) {
	// 1. 生成 Token
	token := util.NewV4().String()

	// 2. 获取 BotUserID（这里简化处理，实际应该从应用配置中获取）
	botUserID := "bot_" + req.AppID

	// 3. 创建 Incoming Webhook 记录
	webhook := models.OpenIncomingWebhook{
		Token:     token,
		AppID:     req.AppID,
		GroupID:   req.GroupID,
		BotUserID: botUserID,
		Name:      req.Name,
		Status:    1,
	}

	if err := l.svcCtx.DB.Create(&webhook).Error; err != nil {
		return nil, errors.New("创建 Incoming Webhook 失败")
	}

	// 4. 构造完整的 Webhook URL
	webhookURL := fmt.Sprintf("/api/open/v1/webhook/incoming/%s", token)

	// 5. 返回结果
	return &types.CreateIncomingWebhookRes{
		Webhook: types.IncomingWebhookInfo{
			ID:         fmt.Sprintf("%d", webhook.ID),
			Token:      webhook.Token,
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
