package logic

import (
	"context"

	"beaver/app/friend/friend_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFriendLogic {
	return &IsFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFriendLogic) IsFriend(in *user_rpc.IsFriendReq) (*user_rpc.IsFriendRes, error) {
	var friend friend_models.FriendModel
	isFriend := friend.IsFriend(l.svcCtx.DB, in.User1, in.User2)

	return &user_rpc.IsFriendRes{
		IsFriend: isFriend,
	}, nil
}
