package logic

import (
	"context"

	"beaver/app/call/call_rpc/internal/svc"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserStatusLogic {
	return &GetUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 供 Chat/User 等服务查询用户是否忙碌
func (l *GetUserStatusLogic) GetUserStatus(in *call_rpc.GetUserStatusReq) (*call_rpc.GetUserStatusRes, error) {
	// todo: add your logic here and delete this line

	return &call_rpc.GetUserStatusRes{}, nil
}
