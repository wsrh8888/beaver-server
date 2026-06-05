package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

const friendVerifyActionDelete int32 = 1

type DeleteFriendVerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFriendVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendVerifyLogic {
	return &DeleteFriendVerifyLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// DeleteFriendVerify 管理后台：删除单条好友验证记录。
// admin 职责：校验 verifyId，映射为 UpdateFriendVerifies 删除 action。
// RPC 职责：UpdateFriendVerifies 统一处理验证记录删除。
func (l *DeleteFriendVerifyLogic) DeleteFriendVerify(req *types.DeleteFriendVerifyReq) (resp *types.DeleteFriendVerifyRes, err error) {
	if req.VerifyID == "" {
		return nil, errors.New("验证记录ID不能为空")
	}

	_, err = l.svcCtx.FriendRpc.UpdateFriendVerifies(l.ctx, &friend_rpc.UpdateFriendVerifiesReq{
		VerifyIds: []string{req.VerifyID},
		Action:    friendVerifyActionDelete,
	})
	if err != nil {
		l.Errorf("删除好友验证失败: %v", err)
		return nil, err
	}
	return &types.DeleteFriendVerifyRes{}, nil
}
