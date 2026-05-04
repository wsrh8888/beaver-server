package message

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendRichTextMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送富文本消息（Markdown）
func NewSendRichTextMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendRichTextMessageLogic {
	return &SendRichTextMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendRichTextMessageLogic) SendRichTextMessage(req *types.SendRichTextMessageReq) (resp *types.SendRichTextMessageRes, err error) {
	// todo: add your logic here and delete this line

	return
}
