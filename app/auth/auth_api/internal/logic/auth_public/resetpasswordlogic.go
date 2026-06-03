package auth_public

import (
	"context"
	"fmt"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/auth/auth_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) (*types.ResetPasswordRes, error) {
	searchRes, err := l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
		Keyword: req.Email,
		Type:    "email",
	})
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	if err = l.verifyEmailCode(req.Email, req.Code, "reset_password"); err != nil {
		return nil, err
	}

	var credential auth_models.AuthCredentialModel
	if err := l.svcCtx.DB.Take(&credential, "user_id = ?", searchRes.UserInfo.UserId).Error; err != nil {
		return nil, fmt.Errorf("重置密码失败")
	}
	credential.Password = pwd.HahPwd(req.Password)
	if err := l.svcCtx.DB.Save(&credential).Error; err != nil {
		return nil, fmt.Errorf("重置密码失败")
	}

	logx.Infof("用户 %s 密码重置成功", req.Email)
	return &types.ResetPasswordRes{}, nil
}

func (l *ResetPasswordLogic) verifyEmailCode(email, code, codeType string) error {
	codeKey := fmt.Sprintf("email_code_%s_%s", email, codeType)
	storedCode, err := l.svcCtx.Redis.Get(codeKey).Result()
	if err != nil {
		return fmt.Errorf("验证码已过期或不存在")
	}
	if storedCode != code {
		return fmt.Errorf("验证码错误")
	}
	l.svcCtx.Redis.Del(codeKey)
	return nil
}
