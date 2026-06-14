package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// DeleteUser 管理后台：删除用户（软删）。
func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (resp *types.DeleteUserRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	_, err = l.svcCtx.UserRpc.DeleteUsers(l.ctx, &user_rpc.DeleteUsersReq{
		UserIds: []string{req.UserID},
	})
	if err != nil {
		l.Errorf("删除用户失败: %v", err)
		return nil, err
	}
	return &types.DeleteUserRes{}, nil
}
