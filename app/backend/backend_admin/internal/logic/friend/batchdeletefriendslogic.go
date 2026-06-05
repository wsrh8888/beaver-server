package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteFriendsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchDeleteFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteFriendsLogic {
	return &BatchDeleteFriendsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// BatchDeleteFriends 管理后台：批量强制删除好友关系。
// admin 职责：校验 ids 非空，映射为 UpdateFriends 批量物理删除。
// RPC 职责：UpdateFriends 统一处理，可复用于其他服务的批量运维能力。
func (l *BatchDeleteFriendsLogic) BatchDeleteFriends(req *types.BatchDeleteFriendsReq) (resp *types.BatchDeleteFriendsRes, err error) {
	if len(req.Ids) == 0 {
		return nil, errors.New("请选择要删除的好友关系")
	}

	_, err = l.svcCtx.FriendRpc.UpdateFriends(l.ctx, &friend_rpc.UpdateFriendsReq{
		RelationIds: req.Ids,
		Action:      friendActionHardDelete,
	})
	if err != nil {
		l.Errorf("批量删除好友失败: %v", err)
		return nil, err
	}
	return &types.BatchDeleteFriendsRes{}, nil
}
