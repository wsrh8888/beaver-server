package bot

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BotStreamMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Bot 流式发送消息（SSE）
func NewBotStreamMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BotStreamMessageLogic {
	return &BotStreamMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BotStreamMessageLogic) BotStreamMessage(req *types.BotStreamMessageReq) (resp *types.BotStreamChunk, err error) {
	// TODO: 流式消息发送，需要 WebSocket/SSE 支持
	appID, _ := l.ctx.Value("appID").(string)
	logx.Infof("Bot 流式发送消息: appID=%s, conversationID=%s", appID, req.ConversationID)

	return nil, errors.New("流式消息功能暂未实现")
}
