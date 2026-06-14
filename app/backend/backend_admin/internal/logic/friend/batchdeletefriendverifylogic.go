package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteFriendVerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchDeleteFriendVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteFriendVerifyLogic {
	return &BatchDeleteFriendVerifyLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// BatchDeleteFriendVerify 管理后台：批量删除好友验证记录。
// admin 职责：校验 ids 非空，映射为 UpdateFriendVerifies 批量删除。
// RPC 职责：UpdateFriendVerifies 统一处理，可复用于其他服务。
func (l *BatchDeleteFriendVerifyLogic) BatchDeleteFriendVerify(req *types.BatchDeleteFriendVerifyReq) (resp *types.BatchDeleteFriendVerifyRes, err error) {
	if len(req.Ids) == 0 {
		return nil, errors.New("请选择要删除的验证记录")
	}

	_, err = l.svcCtx.FriendRpc.UpdateFriendVerifies(l.ctx, &friend_rpc.UpdateFriendVerifiesReq{
		VerifyIds: req.Ids,
		Action:    friendVerifyActionDelete,
	})
	if err != nil {
		l.Errorf("批量删除好友验证失败: %v", err)
		return nil, err
	}
	return &types.BatchDeleteFriendVerifyRes{}, nil
}
