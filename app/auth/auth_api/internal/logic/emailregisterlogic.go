package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmailRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailRegisterLogic {
	return &EmailRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailRegisterLogic) EmailRegister(req *types.EmailRegisterReq) (resp *types.EmailRegisterRes, err error) {
	// 验证邮箱验证码
	err = l.verifyEmailCode(req.Email, req.Code, "register")
	if err != nil {
		return nil, err
	}

	// 检查用户是否已存在
	var user user_models.UserModel
	err = l.svcCtx.DB.Take(&user, "email = ?", req.Email).Error
	if err == nil {
		return nil, errors.New("该邮箱已被注册")
	}

	// 生成随机昵称
	nickName := fmt.Sprintf("用户%s", req.Email[:strings.Index(req.Email, "@")])

	// 创建用户
	_, err = l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateReq{
		Email:    req.Email,
		NickName: nickName,
		Password: req.Password,
		Source:   2,  // 2: 邮箱注册
		Phone:    "", // 邮箱注册时不传手机号
	})
	if err != nil {
		logx.Errorf("创建用户失败: %v", err)
		return nil, errors.New("注册失败")
	}

	return &types.EmailRegisterRes{
		Message: "注册成功",
	}, nil
}

// 验证邮箱验证码
func (l *EmailRegisterLogic) verifyEmailCode(email, code, codeType string) error {
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
