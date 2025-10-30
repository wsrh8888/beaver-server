package logic

import (
	"context"
	"time"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendVerifySyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友验证数据同步
func NewFriendVerifySyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendVerifySyncLogic {
	return &FriendVerifySyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendVerifySyncLogic) FriendVerifySync(req *types.FriendVerifySyncReq) (resp *types.FriendVerifySyncRes, err error) {
	var friendVerifies []friend_models.FriendVerifyModel

	// 设置默认限制
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	// 查询数据 - 查询当前用户作为发送方或接收方的好友验证记录
	err = l.svcCtx.DB.Where("(send_user_id = ? OR rev_user_id = ?) AND version > ? AND version <= ?",
		req.UserID, req.UserID, req.FromVersion, req.ToVersion).
		Order("version ASC").
		Limit(limit + 1).
		Find(&friendVerifies).Error
	if err != nil {
		l.Errorf("查询好友验证数据失败: %v", err)
		return nil, err
	}

	// 判断是否还有更多数据
	hasMore := len(friendVerifies) > limit
	if hasMore {
		friendVerifies = friendVerifies[:limit]
	}

	// 转换为响应格式
	var friendVerifyItems []types.FriendVerifySyncItem
	var nextVersion int64 = req.FromVersion

	for _, friendVerify := range friendVerifies {
		friendVerifyItems = append(friendVerifyItems, types.FriendVerifySyncItem{
			UUID:       friendVerify.UUID, // 添加UUID字段
			SendUserID: friendVerify.SendUserID,
			RevUserID:  friendVerify.RevUserID,
			SendStatus: friendVerify.SendStatus,
			RevStatus:  friendVerify.RevStatus,
			Message:    friendVerify.Message,
			Source:     friendVerify.Source,
			Version:    friendVerify.Version,
			CreateAt:   time.Time(friendVerify.CreatedAt).Unix(),
			UpdateAt:   time.Time(friendVerify.UpdatedAt).Unix(),
		})

		nextVersion = friendVerify.Version
	}

	// 如果没有更多数据，nextVersion应该是toVersion+1
	if !hasMore {
		nextVersion = req.ToVersion + 1
	}

	resp = &types.FriendVerifySyncRes{
		FriendVerifies: friendVerifyItems,
		HasMore:        hasMore,
		NextVersion:    nextVersion,
	}

	l.Infof("好友验证数据同步完成，用户ID: %s, 返回好友验证记录数: %d, 还有更多: %v", req.UserID, len(friendVerifyItems), hasMore)
	return resp, nil
}
