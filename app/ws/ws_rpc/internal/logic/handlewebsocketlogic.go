package logic

import (
	"context"

	"beaver/app/ws/ws_rpc/internal/svc"
	"beaver/app/ws/ws_rpc/types/ws_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleWebSocketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleWebSocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleWebSocketLogic {
	return &HandleWebSocketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HandleWebSocketLogic) HandleWebSocket(in *ws_rpc.HandleWebSocketRequest) (*ws_rpc.HandleWebSocketResponse, error) {
	// todo: add your logic here and delete this line

	return &ws_rpc.HandleWebSocketResponse{}, nil
}
