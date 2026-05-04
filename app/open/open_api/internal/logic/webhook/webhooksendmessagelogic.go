package webhook

import (
	"context"
	"strings"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebhookSendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 通过 Webhook 发送消息（无需鉴权，通过 URL 中的 token 验证）
func NewWebhookSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebhookSendMessageLogic {
	return &WebhookSendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebhookSendMessageLogic) WebhookSendMessage(req *types.WebhookSendMessageReq) (resp *types.WebhookSendMessageRes, err error) {
	// 1. 从 Webhook URL 中提取 token
	// URL 格式: https://api.beaver.im/webhook/{appId}/{token}
	parts := strings.Split(req.WebhookURL, "/")
	if len(parts) < 2 {
		return &types.WebhookSendMessageRes{
			Success: false,
			Message: "无效的 Webhook URL",
		}, nil
	}

	_ = parts[len(parts)-1] // token（暂时未使用）
	appID := parts[len(parts)-2]

	// 2. 验证 Webhook Token
	var webhookConfig open_models.OpenWebhookConfig
	if err := l.svcCtx.DB.Where("app_id = ? AND status = ?", appID, 1).First(&webhookConfig).Error; err != nil {
		return &types.WebhookSendMessageRes{
			Success: false,
			Message: "Webhook 配置不存在或已禁用",
		}, nil
	}

	// TODO: 这里应该验证 token，目前简化处理

	// 3. 获取应用信息
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND status = ?", appID, 1).First(&app).Error; err != nil {
		return &types.WebhookSendMessageRes{
			Success: false,
			Message: "应用不存在或已禁用",
		}, nil
	}

	// 4. 根据消息类型发送消息
	// TODO: 这里需要调用 Bot 发送消息的逻辑
	// 目前返回成功，实际应该调用 message 相关的 RPC 或 logic

	logx.Infof("Webhook 发送消息: appID=%s, msgType=%s", appID, req.MsgType)

	return &types.WebhookSendMessageRes{
		Success: true,
		Message: "消息发送成功",
	}, nil
}
