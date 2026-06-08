package auth_public

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/auth/auth_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/authlock"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)


type EmailRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewEmailRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailRegisterLogic {
	return &EmailRegisterLogic{
		ctx:    ctx,
		logger: logger.New("email_register"),
		svcCtx: svcCtx,
	}
}

func (l *EmailRegisterLogic) EmailRegister(req *types.EmailRegisterReq) (*types.EmailRegisterRes, error) {
	if err := l.verifyEmailCode(req.Email, req.Code, "register"); err != nil {
		return nil, err
	}

	if _, err := l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
		Keyword: req.Email,
		Type:    "email",
	}); err == nil {
		return nil, errors.New("该邮箱已被注册")
	}

	nickName := fmt.Sprintf("用户%s", req.Email[:strings.Index(req.Email, "@")])
	createRes, err := l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateReq{
		Email: req.Email, NickName: nickName, Source: 2,
	})
	if err != nil {
		logx.Errorf("创建用户失败: %v", err)
		return nil, errors.New("注册失败")
	}

	if err := l.svcCtx.DB.Create(&auth_models.AuthCredentialModel{
		UserID: createRes.UserID, Password: pwd.HahPwd(req.Password),
	}).Error; err != nil {
		logx.Errorf("创建用户凭证失败: %v", err)
		return nil, errors.New("创建用户凭证失败")
	}

	logx.Infof("用户注册成功: userID=%s, email=%s", createRes.UserID, req.Email)
	l.logger.Info(model.LogMsg{
		Text: "邮箱注册成功",
		Data: map[string]interface{}{"userId": createRes.UserID},
	})
	return &types.EmailRegisterRes{}, nil
}

func (l *EmailRegisterLogic) verifyEmailCode(email, code, codeType string) error {
	codeKey := fmt.Sprintf("email_code_%s_%s", email, codeType)
	return authlock.VerifyStoredCode(l.ctx, l.svcCtx.Redis, codeKey, codeType, email, code)
}
