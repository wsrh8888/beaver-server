package message

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCardMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送交互式卡片消息
func NewSendCardMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCardMessageLogic {
	return &SendCardMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendCardMessageLogic) SendCardMessage(req *types.SendCardMessageReq) (resp *types.SendCardMessageRes, err error) {
	// todo: add your logic here and delete this line

	return
}
