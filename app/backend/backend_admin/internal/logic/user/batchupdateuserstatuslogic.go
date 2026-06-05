package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

const userActionBatchStatus int32 = 3

type BatchUpdateUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateUserStatusLogic {
	return &BatchUpdateUserStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// BatchUpdateUserStatus 管理后台：批量更新用户状态。
// admin 职责：校验 ids 与 status 合法性，映射为 UpdateUsers 批量改状态 action。
// RPC 职责：UpdateUsers 统一处理状态变更。
func (l *BatchUpdateUserStatusLogic) BatchUpdateUserStatus(req *types.BatchUpdateUserStatusReq) (resp *types.BatchUpdateUserStatusRes, err error) {
	if len(req.Ids) == 0 {
		return nil, errors.New("请选择要操作的用户")
	}
	if req.Status < 1 || req.Status > 3 {
		return nil, errors.New("无效的状态值")
	}

	status := int32(req.Status)
	_, err = l.svcCtx.UserRpc.UpdateUsers(l.ctx, &user_rpc.UpdateUsersReq{
		UserIds:     req.Ids,
		Action:      userActionBatchStatus,
		PatchStatus: &status,
	})
	if err != nil {
		l.Errorf("批量更新用户状态失败: %v", err)
		return nil, err
	}
	return &types.BatchUpdateUserStatusRes{}, nil
}
