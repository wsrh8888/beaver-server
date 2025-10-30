package logic

import (
	"context"
	"time"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友数据同步
func NewFriendSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendSyncLogic {
	return &FriendSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendSyncLogic) FriendSync(req *types.FriendSyncReq) (resp *types.FriendSyncRes, err error) {
	var friends []friend_models.FriendModel

	// 设置默认限制
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	// 查询数据 - 查询当前用户作为发送方或接收方的好友关系
	err = l.svcCtx.DB.Where("(send_user_id = ? OR rev_user_id = ?) AND version > ? AND version <= ?",
		req.UserID, req.UserID, req.FromVersion, req.ToVersion).
		Order("version ASC").
		Limit(limit + 1).
		Find(&friends).Error
	if err != nil {
		l.Errorf("查询好友数据失败: %v", err)
		return nil, err
	}

	// 判断是否还有更多数据
	hasMore := len(friends) > limit
	if hasMore {
		friends = friends[:limit]
	}

	// 转换为响应格式
	var friendItems []types.FriendSyncItem
	var nextVersion int64 = req.FromVersion

	for _, friend := range friends {
		friendItems = append(friendItems, types.FriendSyncItem{
			UUID:           friend.UUID, // 添加UUID字段
			SendUserID:     friend.SendUserID,
			RevUserID:      friend.RevUserID,
			SendUserNotice: friend.SendUserNotice,
			RevUserNotice:  friend.RevUserNotice,
			IsDeleted:      friend.IsDeleted,
			Version:        friend.Version,
			CreateAt:       time.Time(friend.CreatedAt).Unix(),
			UpdateAt:       time.Time(friend.UpdatedAt).Unix(),
			Source:         friend.Source,
		})

		nextVersion = friend.Version
	}

	// 如果没有更多数据，nextVersion应该是toVersion+1
	if !hasMore {
		nextVersion = req.ToVersion + 1
	}

	resp = &types.FriendSyncRes{
		Friends:     friendItems,
		HasMore:     hasMore,
		NextVersion: nextVersion,
	}

	l.Infof("好友数据同步完成，用户ID: %s, 返回好友关系数: %d, 还有更多: %v", req.UserID, len(friendItems), hasMore)
	return resp, nil
}
