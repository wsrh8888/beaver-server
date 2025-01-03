package logic

import (
	"context"

	"beaver/app/ws/ws_rpc/internal/svc"
	"beaver/app/ws/ws_rpc/types/ws_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendProxyMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendProxyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendProxyMessageLogic {
	return &SendProxyMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendProxyMessageLogic) SendProxyMessage(in *ws_rpc.SendProxyMessageRequest) (*ws_rpc.SendProxyMessageResponse, error) {
	// todo: add your logic here and delete this line

	return &ws_rpc.SendProxyMessageResponse{}, nil
}
