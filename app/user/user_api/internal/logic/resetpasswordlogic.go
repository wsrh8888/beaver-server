package logic

import (
	"context"
	"fmt"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) (resp *types.ResetPasswordRes, err error) {
	// 根据邮箱查找用户
	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	// 验证邮箱验证码
	err = l.verifyEmailCode(req.Email, req.Code, "reset_password")
	if err != nil {
		return nil, err
	}

	// 加密新密码
	hashedPassword := pwd.HahPwd(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败")
	}

	// 更新用户密码
	err = l.svcCtx.DB.Model(&user).Update("password", hashedPassword).Error
	if err != nil {
		return nil, err
	}

	logx.Infof("用户 %s 密码重置成功", req.Email)

	return &types.ResetPasswordRes{}, nil
}

// 验证邮箱验证码
func (l *ResetPasswordLogic) verifyEmailCode(email, code, codeType string) error {
	// 从Redis获取存储的验证码
	codeKey := fmt.Sprintf("email_code_%s_%s", email, codeType)
	storedCode, err := l.svcCtx.Redis.Get(codeKey).Result()
	if err != nil {
		return fmt.Errorf("验证码已过期或不存在")
	}

	// 验证验证码
	if storedCode != code {
		return fmt.Errorf("验证码错误")
	}

	// 验证成功后删除验证码
	l.svcCtx.Redis.Del(codeKey)

	return nil
}
