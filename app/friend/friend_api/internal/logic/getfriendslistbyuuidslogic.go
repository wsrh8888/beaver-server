package logic

import (
	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendsListByUuidsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取好友数据（通过ID）
func NewGetFriendsListByUuidsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendsListByUuidsLogic {
	return &GetFriendsListByUuidsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendsListByUuidsLogic) GetFriendsListByUuids(req *types.GetFriendsListByUuidsReq) (resp *types.GetFriendsListByUuidsRes, err error) {
	if len(req.FriendIds) == 0 {
		return &types.GetFriendsListByUuidsRes{
			Friends: []types.FriendByUuid{},
		}, nil
	}

	// 查询指定ID列表中的好友信息
	var friends []friend_models.FriendModel
	err = l.svcCtx.DB.Where("friend_id IN (?)", req.FriendIds).Find(&friends).Error
	if err != nil {
		l.Errorf("查询好友信息失败: ids=%v, error=%v", req.FriendIds, err)
		return nil, err
	}

	l.Infof("查询到 %d 个好友信息", len(friends))

	// 转换为响应格式
	var friendsList []types.FriendByUuid
	for _, friend := range friends {
		l.Infof("处理好友: ID=%s, SendUserID=%s, RevUserID=%s", friend.FriendID, friend.SendUserID, friend.RevUserID)
		friendsList = append(friendsList, types.FriendByUuid{
			FriendID:       friend.FriendID,
			SendUserID:     friend.SendUserID,
			RevUserID:      friend.RevUserID,
			SendUserNotice: friend.SendUserNotice,
			RevUserNotice:  friend.RevUserNotice,
			Source:         friend.Source,
			IsDeleted:      friend.IsDeleted,
			Version:        friend.Version,
			CreateAt:       time.Time(friend.CreatedAt).UnixMilli(),
			UpdateAt:       time.Time(friend.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetFriendsListByUuidsRes{
		Friends: friendsList,
	}, nil
}
