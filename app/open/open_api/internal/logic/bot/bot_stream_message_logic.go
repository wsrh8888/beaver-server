// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package bot

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
