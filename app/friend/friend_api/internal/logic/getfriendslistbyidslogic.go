package logic

import (
	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendsListByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取好友数据
func NewGetFriendsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendsListByIdsLogic {
	return &GetFriendsListByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendsListByIdsLogic) GetFriendsListByIds(req *types.GetFriendsListByIdsReq) (resp *types.GetFriendsListByIdsRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.Errorf("用户ID为空")
		return nil, errors.New("用户ID不能为空")
	}

	if len(req.FriendIds) == 0 {
		return &types.GetFriendsListByIdsRes{
			Friends: []types.FriendSyncItem{},
		}, nil
	}

	// 查询指定好友ID的好友数据
	var friends []friend_models.FriendModel
	err = l.svcCtx.DB.Where("uuid IN (?) AND (send_user_id = ? OR rev_user_id = ?)", req.FriendIds, userId, userId).Find(&friends).Error
	if err != nil {
		l.Errorf("查询好友数据失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var friendItems []types.FriendSyncItem
	for _, friend := range friends {
		friendItems = append(friendItems, types.FriendSyncItem{
			UUID:           friend.UUID,
			SendUserID:     friend.SendUserID,
			RevUserID:      friend.RevUserID,
			SendUserNotice: friend.SendUserNotice,
			RevUserNotice:  friend.RevUserNotice,
			IsDeleted:      friend.IsDeleted,
			Version:        friend.Version,
			CreateAt:       time.Time(friend.CreatedAt).Unix(),
			Source:         friend.Source,
			UpdateAt:       time.Time(friend.UpdatedAt).Unix(),
		})
	}

	return &types.GetFriendsListByIdsRes{
		Friends: friendItems,
	}, nil
}
