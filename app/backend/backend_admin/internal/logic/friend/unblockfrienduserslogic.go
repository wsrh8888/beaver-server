package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnblockFriendUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnblockFriendUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnblockFriendUsersLogic {
	return &UnblockFriendUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnblockFriendUsersLogic) UnblockFriendUsers(req *types.UnblockFriendUsersReq) (resp *types.UnblockFriendUsersRes, err error) {
	if len(req.Ids) == 0 {
		return nil, errors.New("请选择要解除的黑名单记录")
	}

	_, err = l.svcCtx.FriendRpc.UpdateFriendBlocks(l.ctx, &friend_rpc.UpdateFriendBlocksReq{
		BlockIds: req.Ids,
		Action:   1,
	})
	if err != nil {
		l.Errorf("解除黑名单失败: %v", err)
		return nil, err
	}

	return &types.UnblockFriendUsersRes{}, nil
}
