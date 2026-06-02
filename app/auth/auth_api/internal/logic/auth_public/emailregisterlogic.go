package auth_public

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	auth_models "beaver/app/auth/auth_models"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/pwd"

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

func (l *EmailRegisterLogic) EmailRegister(req *types.EmailRegisterReq) (*types.EmailRegisterRes, error) {
	if err := l.verifyEmailCode(req.Email, req.Code, "register"); err != nil {
		return nil, err
	}

	var user user_models.UserModel
	if err := l.svcCtx.DB.Take(&user, "email = ?", req.Email).Error; err == nil {
		return nil, errors.New("该邮箱已被注册")
	}

	nickName := fmt.Sprintf("用户%s", req.Email[:strings.Index(req.Email, "@")])
	createRes, err := l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateReq{
		Email:    req.Email,
		NickName: nickName,
		Source:   2,
		Phone:    "",
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

	logx.Infof("用户注册成功: userID=%s, email=%s", userID, req.Email)
	return &types.EmailRegisterRes{}, nil
}

func (l *EmailRegisterLogic) verifyEmailCode(email, code, codeType string) error {
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