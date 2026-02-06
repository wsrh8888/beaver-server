package logic

import (
	"context"

	"beaver/app/call/call_rpc/internal/svc"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 供 Call-Api 调用，创建通话记录
func (l *CreateSessionLogic) CreateSession(in *call_rpc.CreateSessionReq) (*call_rpc.CreateSessionRes, error) {
	// todo: add your logic here and delete this line

	return &call_rpc.CreateSessionRes{}, nil
}
