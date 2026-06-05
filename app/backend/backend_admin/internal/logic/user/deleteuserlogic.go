package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

const userActionSoftDelete int32 = 2

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// DeleteUser 管理后台：删除用户（软删）。
// admin 职责：校验 userId，映射为 UpdateUsers 软删除 action。
// RPC 职责：UpdateUsers 统一处理用户状态变更。
func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (resp *types.DeleteUserRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	_, err = l.svcCtx.UserRpc.UpdateUsers(l.ctx, &user_rpc.UpdateUsersReq{
		UserIds: []string{req.UserID},
		Action:  userActionSoftDelete,
	})
	if err != nil {
		l.Errorf("删除用户失败: %v", err)
		return nil, err
	}
	return &types.DeleteUserRes{}, nil
}
