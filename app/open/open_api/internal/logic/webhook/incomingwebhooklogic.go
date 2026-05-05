// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package webhook

import (
	"context"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type IncomingWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 接收外部 Webhook（用于 Jenkins/GitHub 等集成，无需鉴权）
func NewIncomingWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IncomingWebhookLogic {
	return &IncomingWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IncomingWebhookLogic) IncomingWebhook(req *types.IncomingWebhookReq) (resp *types.IncomingWebhookRes, err error) {
	// 1. 验证 Token
	var webhook open_models.OpenIncomingWebhook
	if err := l.svcCtx.DB.Where("token = ? AND status = ?", req.Token, 1).First(&webhook).Error; err != nil {
		return &types.IncomingWebhookRes{
			Success: false,
			Message: "无效的 Webhook Token",
		}, nil
	}

	// 2. 构造消息内容
	content := ""
	if req.MsgType == "text" {
		if text, ok := req.Content["text"].(string); ok {
			content = text
		} else {
			return &types.IncomingWebhookRes{
				Success: false,
				Message: "消息格式错误",
			}, nil
		}
	} else if req.MsgType == "markdown" {
		if md, ok := req.Content["content"].(string); ok {
			content = md
		} else {
			return &types.IncomingWebhookRes{
				Success: false,
				Message: "消息格式错误",
			}, nil
		}
	} else {
		return &types.IncomingWebhookRes{
			Success: false,
			Message: "不支持的消息类型",
		}, nil
	}

	// 3. 构造消息
	var msg *chat_rpc.Msg
	if req.MsgType == "text" {
		msg = &chat_rpc.Msg{
			Type: 1, // 文本消息
			TextMsg: &chat_rpc.TextMsg{
				Content: content,
			},
		}
	} else if req.MsgType == "markdown" {
		// Markdown 也作为文本消息发送
		msg = &chat_rpc.Msg{
			Type: 1,
			TextMsg: &chat_rpc.TextMsg{
				Content: content,
			},
		}
	}

	// 4. 调用 Chat RPC 发送消息到群组
	sendReq := &chat_rpc.SendMsgReq{
		UserId:         webhook.BotUserID,
		ConversationId: webhook.GroupID,
		Msg:            msg,
	}

	_, err = l.svcCtx.ChatRpc.SendMsg(l.ctx, sendReq)
	if err != nil {
		logx.Errorf("发送消息失败: %v", err)
		return &types.IncomingWebhookRes{
			Success: false,
			Message: "发送消息失败",
		}, nil
	}

	logx.Infof("Incoming Webhook 消息发送成功: group_id=%s, bot_id=%s", webhook.GroupID, webhook.BotUserID)

	return &types.IncomingWebhookRes{
		Success: true,
	}, nil
}
