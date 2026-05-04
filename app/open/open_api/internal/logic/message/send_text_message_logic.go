// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package message

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
