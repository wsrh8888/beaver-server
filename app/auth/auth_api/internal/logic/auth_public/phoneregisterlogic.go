package auth_public

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	auth_models "beaver/app/auth/auth_models"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/pwd"

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

func (l *PhoneRegisterLogic) PhoneRegister(req *types.PhoneRegisterReq) (*types.PhoneRegisterRes, error) {
	if err := l.verifyPhoneCode(req.Phone, req.Code, "register"); err != nil {
		return nil, err
	}

	var user user_models.UserModel
	if err := l.svcCtx.DB.Take(&user, "phone = ?", req.Phone).Error; err == nil {
		return nil, errors.New("该手机号已被注册")
	}

	createRes, err := l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateReq{
		Phone:  req.Phone,
		Source: 1,
	})
	if err != nil {
		logx.Errorf("创建用户失败: %v", err)
		return nil, errors.New("注册失败")
	}

	userID := createRes.UserID
	hashedPassword := pwd.HahPwd(req.Password)
	credential := auth_models.AuthCredentialModel{
		UserID:   userID,
		Password: hashedPassword,
	}

	if err := l.svcCtx.DB.Create(&credential).Error; err != nil {
		logx.Errorf("创建用户凭证失败: %v", err)
		l.svcCtx.DB.Where("user_id = ?", userID).Delete(&user_models.UserModel{})
		return nil, errors.New("创建用户凭证失败")
	}

	logx.Infof("用户注册成功: userID=%s, phone=%s", userID, req.Phone)
	return &types.PhoneRegisterRes{}, nil
}

func (l *PhoneRegisterLogic) verifyPhoneCode(phone, code, codeType string) error {
	codeKey := fmt.Sprintf("phone_code_%s_%s", phone, codeType)
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