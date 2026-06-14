package oauth_secret

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserInfoByQuickLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoByQuickLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByQuickLoginLogic {
	return &GetUserInfoByQuickLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoByQuickLoginLogic) GetUserInfoByQuickLogin(req *types.GetUserInfoByQuickLoginReq) (resp *types.GetUserInfoByQuickLoginRes, err error) {
	if _, err := verifyApp(l.svcCtx.DB, req.AppID, req.AppSecret); err != nil {
		return nil, err
	}

	oauthCode, err := findOAuthCode(l.svcCtx.DB, req.AppID, req.AuthCode)
	if err != nil {
		return nil, err
	}
	if oauthCode.Scene != "pc_scan" && oauthCode.Scene != "h5_sso" {
		return nil, errors.New("授权码场景无效")
	}

	if err := l.svcCtx.DB.Model(oauthCode).Update("used", true).Error; err != nil {
		logx.Errorf("标记 authCode 已使用失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	userInfoRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: oauthCode.UserID,
	})
	if err != nil {
		logx.Errorf("查询用户信息失败: %v", err)
		return nil, errors.New("获取用户信息失败")
	}
	if userInfoRes.UserInfo == nil {
		return nil, errors.New("用户不存在")
	}

	return &types.GetUserInfoByQuickLoginRes{
		UserID:   userInfoRes.UserInfo.UserId,
		NickName: userInfoRes.UserInfo.NickName,
		Avatar:   userInfoRes.UserInfo.Avatar,
		Phone:    "",
		Email:    userInfoRes.UserInfo.Email,
	}, nil
}

func verifyApp(db *gorm.DB, appID, appSecret string) (*open_models.OpenApp, error) {
	if appID == "" || appSecret == "" {
		return nil, errors.New("应用凭证不完整")
	}
	var app open_models.OpenApp
	if err := db.Where("app_id = ? AND app_secret = ? AND status = ?", appID, appSecret, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或凭证错误")
	}
	return &app, nil
}

func findOAuthCode(db *gorm.DB, appID, code string) (*open_models.OpenOAuthCode, error) {
	var record open_models.OpenOAuthCode
	if err := db.Where("code = ? AND app_id = ?", code, appID).First(&record).Error; err != nil {
		return nil, errors.New("授权码无效")
	}
	if record.Used {
		return nil, errors.New("授权码已使用")
	}
	if time.Now().Unix() > record.ExpiresAt {
		return nil, errors.New("授权码已过期")
	}
	return &record, nil
}
