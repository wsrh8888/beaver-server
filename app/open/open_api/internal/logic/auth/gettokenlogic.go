package auth

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

type GetTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取访问令牌
func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenLogic) GetToken(req *types.GetTokenReq) (resp *types.GetTokenRes, err error) {
	// 1. 验证应用凭证
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND app_secret = ? AND status = ?", req.AppID, req.AppSecret, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用 ID 或密钥错误")
	}

	// 2. 生成 Access Token
	accessTokenBytes := make([]byte, 32)
	_, _ = rand.Read(accessTokenBytes)
	accessToken := hex.EncodeToString(accessTokenBytes)

	// 3. 生成 Refresh Token
	refreshTokenBytes := make([]byte, 32)
	_, _ = rand.Read(refreshTokenBytes)
	refreshToken := hex.EncodeToString(refreshTokenBytes)

	// 4. 保存 Access Token（客户端凭证模式没有用户ID）
	now := time.Now()
	expiresAt := now.Add(2 * time.Hour).Unix() // 2小时过期
	tokenRecord := open_models.OpenAccessToken{
		AppID:        req.AppID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Scope:        app.Scopes,
		UserID:       "", // 客户端凭证模式没有用户
	}

	if err := l.svcCtx.DB.Create(&tokenRecord).Error; err != nil {
		logx.Errorf("创建访问令牌失败: %v", err)
		return nil, errors.New("生成令牌失败")
	}

	return &types.GetTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200, // 2小时
		TokenType:    "Bearer",
	}, nil
}
