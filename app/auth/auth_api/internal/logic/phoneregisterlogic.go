package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PhoneRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPhoneRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PhoneRegisterLogic {
	return &PhoneRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PhoneRegisterLogic) PhoneRegister(req *types.PhoneRegisterReq) (resp *types.PhoneRegisterRes, err error) {
	// 验证手机验证码
	err = l.verifyPhoneCode(req.Phone, req.Code, "register")
	if err != nil {
		return nil, err
	}

	// 检查用户是否已存在
	var user user_models.UserModel
	err = l.svcCtx.DB.Take(&user, "phone = ?", req.Phone).Error
	if err == nil {
		return nil, errors.New("该手机号已被注册")
	}

	// 创建用户
	_, err = l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateReq{
		Phone:    req.Phone,
		Password: req.Password,
		Source:   1, // 1: 手机号注册
	})
	if err != nil {
		logx.Errorf("创建用户失败: %v", err)
		return nil, errors.New("注册失败")
	}

	return &types.PhoneRegisterRes{
		Message: "注册成功",
	}, nil
}

// 验证手机验证码
func (l *PhoneRegisterLogic) verifyPhoneCode(phone, code, codeType string) error {
	// 从Redis获取存储的验证码
	codeKey := fmt.Sprintf("phone_code_%s_%s", phone, codeType)
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
