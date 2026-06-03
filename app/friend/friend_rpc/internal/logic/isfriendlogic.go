package logic

import (
	"context"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

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

func (l *IsFriendLogic) IsFriend(in *friend_rpc.IsFriendReq) (*friend_rpc.IsFriendRes, error) {
	var friend friend_models.FriendModel
	isFriend := friend.IsFriend(l.svcCtx.DB, in.UserA, in.UserB)
	return &friend_rpc.IsFriendRes{IsFriend: isFriend}, nil
}
