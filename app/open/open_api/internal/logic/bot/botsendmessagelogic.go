package bot

import (
	"context"
	"errors"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

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

func (l *BotSendMessageLogic) BotSendMessage(req *types.BotSendMessageReq) (resp *types.BotSendMessageRes, err error) {
	// 1. 根据 AppID 查询 Bot 用户 ID
	// TODO: 从数据库或缓存中获取 Bot 对应的 UserID
	botUserID := "" // 需要根据 req.AppID 查询

	if botUserID == "" {
		return nil, errors.New("Bot 未配置或应用不存在")
	}

	// 2. 构造 Chat RPC 请求
	msg := &chat_rpc.Msg{
		Type: 1, // 文本消息
		TextMsg: &chat_rpc.TextMsg{
			Content: req.Content,
		},
	}

	rpcReq := &chat_rpc.SendMsgReq{
		UserId:         botUserID,
		ConversationId: req.ConversationID,
		MessageId:      req.IdempotentKey,
		Msg:            msg,
	}

	// 3. 调用 Chat RPC 发送消息
	rpcResp, err := l.svcCtx.ChatRpc.SendMsg(l.ctx, rpcReq)
	if err != nil {
		logx.Errorf("调用 Chat RPC 发送消息失败: %v", err)
		return nil, errors.New("发送消息失败")
	}

	return &types.BotSendMessageRes{
		MessageID: rpcResp.MessageId,
	}, nil
}
