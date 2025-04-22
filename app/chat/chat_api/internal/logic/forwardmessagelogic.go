package logic

import (
	"context"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForwardMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForwardMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForwardMessageLogic {
	return &ForwardMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForwardMessageLogic) ForwardMessage(req *types.ForwardMessageReq) (resp *types.ForwardMessageRes, err error) {
	// todo: add your logic here and delete this line

	return
}
