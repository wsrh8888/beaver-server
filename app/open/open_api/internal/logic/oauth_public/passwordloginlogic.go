package oauth_public

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/app/open/open_api/internal/logic/oauthutil"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PasswordLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 账号密码登录
func NewPasswordLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PasswordLoginLogic {
	return &PasswordLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PasswordLoginLogic) PasswordLogin(req *types.PasswordLoginReq) (resp *types.PasswordLoginRes, err error) {
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		logx.Errorf("应用不存在: appId=%s, err=%v", req.AppID, err)
		return nil, fmt.Errorf("应用不存在")
	}
	if app.Status != 1 {
		return nil, fmt.Errorf("应用未启用")
	}

	userRes, err := l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
		Keyword: req.Account,
		Type:    "email",
	})
	if err != nil {
		userRes, err = l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
			Keyword: req.Account,
			Type:    "userId",
		})
	}
	if err != nil {
		logx.Errorf("用户不存在: account=%s, err=%v", req.Account, err)
		return nil, errors.New("账号或密码错误")
	}

	verifyRes, err := l.svcCtx.AuthRpc.VerifyPassword(l.ctx, &auth_rpc.VerifyPasswordReq{
		UserId:   userRes.UserInfo.UserId,
		Password: req.Password,
	})
	if err != nil || !verifyRes.Valid {
		logx.Errorf("密码错误: userId=%s", userRes.UserInfo.UserId)
		return nil, errors.New("账号或密码错误")
	}

	code, expireIn, err := oauthutil.CreateOAuthCode(l.svcCtx.DB, req.AppID, userRes.UserInfo.UserId, "password")
	if err != nil {
		logx.Errorf("生成授权码失败: userId=%s, err=%v", userRes.UserInfo.UserId, err)
		return nil, fmt.Errorf("服务内部异常")
	}

	logx.Infof("账号密码登录成功: userId=%s, appId=%s, account=%s", userRes.UserInfo.UserId, req.AppID, req.Account)

	return &types.PasswordLoginRes{
		Code:     code,
		ExpireIn: expireIn,
	}, nil
}
