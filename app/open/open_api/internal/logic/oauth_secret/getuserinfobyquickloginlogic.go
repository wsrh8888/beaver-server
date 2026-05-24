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
)

type GetUserInfoByQuickLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// PC 端快捷登录（用 authCode 换取用户信息）
func NewGetUserInfoByQuickLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByQuickLoginLogic {
	return &GetUserInfoByQuickLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoByQuickLoginLogic) GetUserInfoByQuickLogin(req *types.GetUserInfoByQuickLoginReq) (resp *types.GetUserInfoByQuickLoginRes, err error) {
	// 1. 验证参数
	if req.AppID == "" || req.AppSecret == "" || req.AuthCode == "" {
		return nil, errors.New("参数不完整")
	}

	// 2. 验证应用合法性（appId + appSecret）
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND app_secret = ?", req.AppID, req.AppSecret).First(&app).Error; err != nil {
		logx.Errorf("应用验证失败: appId=%s, err=%v", req.AppID, err)
		return nil, errors.New("应用不存在或 appSecret 错误")
	}

	// 3. 从 authCode 查询用户信息
	var oauthCode open_models.OpenOAuthCode
	if err := l.svcCtx.DB.Where("code = ? AND app_id = ?", req.AuthCode, req.AppID).First(&oauthCode).Error; err != nil {
		logx.Errorf("authCode 查询失败: authCode=%s, err=%v", req.AuthCode, err)
		return nil, errors.New("授权码无效")
	}

	// 4. 检查 authCode 是否过期
	if time.Now().Unix() > oauthCode.ExpiresAt {
		return nil, errors.New("授权码已过期")
	}

	// 5. 检查 authCode 是否已使用
	if oauthCode.Used {
		return nil, errors.New("授权码已使用")
	}

	// 6. 标记 authCode 为已使用
	oauthCode.Used = true
	l.svcCtx.DB.Save(&oauthCode)

	// 7. 调用 UserRpc 获取用户信息
	userInfoRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: oauthCode.UserID,
	})
	if err != nil {
		logx.Errorf("获取用户信息失败: %v", err)
		return nil, errors.New("获取用户信息失败")
	}

	if userInfoRes.UserInfo == nil {
		return nil, errors.New("用户不存在")
	}

	// 8. 返回用户信息
	return &types.GetUserInfoByQuickLoginRes{
		UserID:   userInfoRes.UserInfo.UserId,
		NickName: userInfoRes.UserInfo.NickName,
		Avatar:   userInfoRes.UserInfo.Avatar,
		Phone:    "",
		Email:    userInfoRes.UserInfo.Email,
	}, nil
}
