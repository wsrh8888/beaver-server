package logic

import (
	"context"
	"errors"

	"beaver/app/auth/auth_models"
	"beaver/app/auth/auth_rpc/internal/svc"
	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(in *auth_rpc.ResetPasswordReq) (*auth_rpc.ResetPasswordRes, error) {
	if in.UserId == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if in.NewPassword == "" {
		return nil, errors.New("新密码不能为空")
	}

	var credential auth_models.AuthCredentialModel
	if err := l.svcCtx.DB.Take(&credential, "user_id = ?", in.UserId).Error; err != nil {
		logx.Errorf("查询用户凭证失败: %v", err)
		return nil, errors.New("用户凭证不存在")
	}

	credential.Password = pwd.HahPwd(in.NewPassword)
	if err := l.svcCtx.DB.Save(&credential).Error; err != nil {
		logx.Errorf("重置密码失败: %v", err)
		return nil, errors.New("重置密码失败")
	}

	logx.Infof("密码重置成功: userID=%s", in.UserId)
	return &auth_rpc.ResetPasswordRes{Success: true}, nil
}
