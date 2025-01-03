package logic

import (
	"context"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendIdsLogic {
	return &GetFriendIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendIdsLogic) GetFriendIds(in *friend_rpc.GetFriendIdsRequest) (*friend_rpc.GetFriendIdsResponse, error) {
	var friends []friend_models.FriendModel
	err := l.svcCtx.DB.Where("(send_user_id = ? OR rev_user_id = ?) AND is_deleted = false", in.UserId, in.UserId).Find(&friends).Error
	if err != nil {
		logx.Errorf("failed to query friends: %v", err)
		return nil, err
	}

	var friendIds []string
	for _, friend := range friends {
		if friend.SendUserId == in.UserId {
			friendIds = append(friendIds, friend.RevUserId)
		} else {
			friendIds = append(friendIds, friend.SendUserId)
		}
	}

	return &friend_rpc.GetFriendIdsResponse{
		FriendIds: friendIds,
	}, nil
}
