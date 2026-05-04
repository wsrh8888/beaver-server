package message

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendTextMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送文本消息
func NewSendTextMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendTextMessageLogic {
	return &SendTextMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendTextMessageLogic) SendTextMessage(req *types.SendTextMessageReq) (resp *types.SendTextMessageRes, err error) {
	// TODO: 需要构造会话 ID 和获取发送者 ID
	// 目前开放平台的消息发送主要通过 Bot API，这个接口保留作为扩展
	return nil, errors.New("请使用 Bot API 发送消息")
}
