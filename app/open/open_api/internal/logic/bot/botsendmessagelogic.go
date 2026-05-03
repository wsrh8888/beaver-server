package bot

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type BotSendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Bot 主动发送消息（对标飞书/钉钉 Bot API）
func NewBotSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BotSendMessageLogic {
	return &BotSendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Bot 主动发送消息（对标飞书/钉钉 Bot API）
func (l *BotSendMessageLogic) BotSendMessage(req *types.BotSendMessageReq) (resp *types.BotSendMessageRes, err error) {
	// 1. 验证 App
	var app open_models.OpenApp
	err = l.svcCtx.DB.Where("app_id = ? AND status = ?", req.AppID, 1).First(&app).Error
	if err != nil {
		return nil, errors.New("应用不存在或已禁用")
	}

	// 2. 验证 Bot 用户
	if app.BotUserID == "" {
		return nil, errors.New("Bot 未配置")
	}

	// 3. 构建消息内容
	msgContent := map[string]interface{}{
		"type":    req.MsgType, // text/markdown/richtext/html
		"content": req.Content,
	}

	if req.Metadata != nil && len(req.Metadata) > 0 {
		msgContent["metadata"] = req.Metadata
	}

	// 4. 调用 chat_rpc 发送消息
	messageID := uuid.New().String()

	// TODO: 构造正确的消息格式
	// 将消息内容序列化为 JSON
	// contentJSON, _ := json.Marshal(msgContent)

	_, err = l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
		UserId:    app.BotUserID, // Bot 作为发送者
		MessageId: messageID,
		// ConversationId: req.ConversationID, // 目标会话（私聊或群聊）
		// Msg:            contentJSON,
	})

	if err != nil {
		return nil, fmt.Errorf("发送消息失败: %w", err)
	}

	return &types.BotSendMessageRes{
		MessageID: messageID,
	}, nil
}
