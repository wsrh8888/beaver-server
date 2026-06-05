package logic

import (
	"context"
	"errors"

	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetUserPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	return &ResetUserPasswordLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// ResetUserPassword 管理后台：重置用户密码。
// admin 职责：校验 userId/newPassword，调用 AuthRpc 重置凭证（密码属认证域，不走 UserRpc）。
// RPC 职责：AuthRpc.ResetPassword 管理密码哈希与凭证状态。
func (l *ResetUserPasswordLogic) ResetUserPassword(req *types.ResetUserPasswordReq) (resp *types.ResetUserPasswordRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if req.NewPassword == "" {
		return nil, errors.New("新密码不能为空")
	}

	_, err = l.svcCtx.AuthRpc.ResetPassword(l.ctx, &auth_rpc.ResetPasswordReq{
		UserId:      req.UserID,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		l.Errorf("重置用户密码失败: %v", err)
		return nil, err
	}
	return &types.ResetUserPasswordRes{}, nil
}
