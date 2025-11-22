package logic

import (
	"context"
	"time"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncFriendsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的好友版本
func NewGetSyncFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncFriendsLogic {
	return &GetSyncFriendsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncFriendsLogic) GetSyncFriends(req *types.GetSyncFriendsReq) (resp *types.GetSyncFriendsRes, err error) {
	// 调用Friend RPC获取好友版本信息
	friendResp, err := l.svcCtx.FriendRpc.GetFriendVersions(l.ctx, &friend_rpc.GetFriendVersionsReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.Errorf("获取好友版本信息失败: userId=%s, since=%d, error=%v", req.UserID, req.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个好友版本信息", len(friendResp.FriendVersions))

	// 转换为响应格式，确保返回空数组而不是null
	friendVersions := make([]types.FriendVersionItem, 0)
	if friendResp.FriendVersions != nil {
		for _, friend := range friendResp.FriendVersions {
			friendVersions = append(friendVersions, types.FriendVersionItem{
				Id:      friend.Id,
				Version: friend.Version,
			})
		}
	}

	return &types.GetSyncFriendsRes{
		FriendVersions:  friendVersions,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
