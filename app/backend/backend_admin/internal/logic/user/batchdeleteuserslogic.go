package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchDeleteUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteUsersLogic {
	return &BatchDeleteUsersLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// BatchDeleteUsers 管理后台：批量软删除用户。
// admin 职责：校验 ids 非空，映射为 UpdateUsers 批量软删除。
// RPC 职责：UpdateUsers 可复用的批量运维能力。
func (l *BatchDeleteUsersLogic) BatchDeleteUsers(req *types.BatchDeleteUsersReq) (resp *types.BatchDeleteUsersRes, err error) {
	if len(req.Ids) == 0 {
		return nil, errors.New("请选择要删除的用户")
	}

	_, err = l.svcCtx.UserRpc.UpdateUsers(l.ctx, &user_rpc.UpdateUsersReq{
		UserIds: req.Ids,
		Action:  userActionSoftDelete,
	})
	if err != nil {
		l.Errorf("批量删除用户失败: %v", err)
		return nil, err
	}
	return &types.BatchDeleteUsersRes{}, nil
}
