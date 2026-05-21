package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExchangeTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeTokenLogic {
	return &ExchangeTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExchangeTokenLogic) ExchangeToken(in *open_rpc.ExchangeTokenReq) (*open_rpc.ExchangeTokenRes, error) {
	// 1. 验证应用是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND status = ?", in.AppId, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或已禁用")
	}

	// 2. 查询授权码
	var authCode open_models.OpenOAuthCode
	if err := l.svcCtx.DB.Where("code = ? AND app_id = ? AND used = ?", in.Code, in.AppId, false).First(&authCode).Error; err != nil {
		return nil, errors.New("授权码无效或已使用")
	}

	// 3. 检查授权码是否过期
	if time.Now().Unix() > authCode.ExpiresAt {
		return nil, errors.New("授权码已过期")
	}

	// 4. 标记授权码为已使用
	if err := l.svcCtx.DB.Model(&authCode).Update("used", true).Error; err != nil {
		logx.Errorf("更新授权码状态失败: %v", err)
		return nil, errors.New("处理授权码失败")
	}

	// 5. 生成 Access Token
	accessTokenBytes := make([]byte, 32)
	_, _ = rand.Read(accessTokenBytes)
	accessToken := hex.EncodeToString(accessTokenBytes)

	// 6. 生成 Refresh Token
	refreshTokenBytes := make([]byte, 32)
	_, _ = rand.Read(refreshTokenBytes)
	refreshToken := hex.EncodeToString(refreshTokenBytes)

	// 7. 保存 Access Token
	now := time.Now()
	expiresAt := now.Add(2 * time.Hour).Unix()                    // access_token 2小时过期
	refreshTokenExpiresAt := now.Add(180 * 24 * time.Hour).Unix() // refresh_token 180天过期
	tokenRecord := open_models.OpenOAuthToken{
		AppID:                 in.AppId,
		Token:                 accessToken,
		RefreshToken:          refreshToken,
		ExpiresAt:             expiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		Scope:                 authCode.Scope,
		UserID:                authCode.UserID,
	}

	if err := l.svcCtx.DB.Create(&tokenRecord).Error; err != nil {
		logx.Errorf("创建访问令牌失败: %v", err)
		return nil, errors.New("生成令牌失败")
	}

	return &open_rpc.ExchangeTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200, // 2小时
		TokenType:    "Bearer",
		UserId:       authCode.UserID,
	}, nil
}
