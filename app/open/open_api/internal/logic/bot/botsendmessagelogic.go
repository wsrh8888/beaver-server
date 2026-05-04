package bot

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
