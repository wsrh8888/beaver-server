package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

const friendActionRestore int32 = 2

type RestoreFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRestoreFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreFriendLogic {
	return &RestoreFriendLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// RestoreFriend 管理后台：恢复软删除的好友关系。
// admin 职责：校验 friendId，映射为 UpdateFriends 恢复 action。
// RPC 职责：UpdateFriends 统一处理关系状态变更。
func (l *RestoreFriendLogic) RestoreFriend(req *types.RestoreFriendReq) (resp *types.RestoreFriendRes, err error) {
	if req.FriendID == "" {
		return nil, errors.New("好友关系ID不能为空")
	}

	_, err = l.svcCtx.FriendRpc.UpdateFriends(l.ctx, &friend_rpc.UpdateFriendsReq{
		RelationIds: []string{req.FriendID},
		Action:      friendActionRestore,
	})
	if err != nil {
		l.Errorf("恢复好友失败: %v", err)
		return nil, err
	}
	return &types.RestoreFriendRes{}, nil
}
