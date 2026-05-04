// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package message

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendImageMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送图片消息
func NewSendImageMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendImageMessageLogic {
	return &SendImageMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendImageMessageLogic) SendImageMessage(req *types.SendImageMessageReq) (resp *types.SendImageMessageRes, err error) {
	// todo: add your logic here and delete this line

	return
}
