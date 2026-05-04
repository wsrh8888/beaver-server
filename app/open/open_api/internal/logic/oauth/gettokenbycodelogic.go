package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenByCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用授权码换取 Token
func NewGetTokenByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenByCodeLogic {
	return &GetTokenByCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenByCodeLogic) GetTokenByCode(req *types.GetTokenByCodeReq) (resp *types.GetTokenByCodeRes, err error) {
	// 1. 验证应用凭证
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND app_secret = ? AND status = ?", req.AppID, req.AppSecret, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用 ID 或密钥错误")
	}

	// 2. 查询授权码
	var authCode open_models.OpenAuthCode
	if err := l.svcCtx.DB.Where("code = ? AND app_id = ? AND used = ?", req.Code, req.AppID, false).First(&authCode).Error; err != nil {
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
	expiresAt := now.Add(2 * time.Hour).Unix() // 2小时过期
	tokenRecord := open_models.OpenAccessToken{
		AppID:        req.AppID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Scope:        authCode.Scope,
		UserID:       authCode.UserID,
	}

	if err := l.svcCtx.DB.Create(&tokenRecord).Error; err != nil {
		logx.Errorf("创建访问令牌失败: %v", err)
		return nil, errors.New("生成令牌失败")
	}

	return &types.GetTokenByCodeRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200, // 2小时
		TokenType:    "Bearer",
	}, nil
}
