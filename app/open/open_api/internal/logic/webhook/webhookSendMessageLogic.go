package webhook

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	models "beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebhookSendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebhookSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebhookSendMessageLogic {
	return &WebhookSendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebhookSendMessageLogic) WebhookSendMessage(req *types.WebhookSendMessageReq) (resp *types.WebhookSendMessageRes, err error) {
	// 1. 从 URL 参数中获取 token
	token := l.ctx.Value("webhook_token")
	if token == nil || token.(string) == "" {
		return nil, errors.New("缺少 token 参数")
	}

	// 2. 验证 Webhook 配置
	var config models.OpenWebhookConfig
	err = l.svcCtx.DB.Where("target_url LIKE ? AND status = ?", "%"+token.(string)+"%", 1).First(&config).Error
	if err != nil {
		return nil, errors.New("Webhook 无效或已禁用")
	}

	// 3. 构建消息内容(暂时不使用,后续调用 chat_rpc 时需要)
	_ = req

	// 4. TODO: 调用 chat_rpc 发送消息到群组
	// 这里需要根据 AppID 找到对应的群聊，然后发送消息
	// 暂时返回成功，实际实现需要调用 chat_rpc

	logx.Infof("Webhook 消息发送成功: app_id=%s, msg_type=%s",
		config.AppID, req.MsgType)

	return &types.WebhookSendMessageRes{
		Success: true,
	}, nil
}
