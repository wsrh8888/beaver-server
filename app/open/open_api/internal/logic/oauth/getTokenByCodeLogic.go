// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package oauth

import (
	"context"
	"errors"
	"time"

	"beaver-server/app/open/open_api/internal/svc"
	"beaver-server/app/open/open_api/internal/types"
	models "beaver-server/app/open/open_models"

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
	// 1. 验证应用
	var app models.OpenApp
	err = l.svcCtx.DB.Where("app_id = ? AND app_secret = ?", req.AppID, req.AppSecret).First(&app).Error
	if err != nil {
		return nil, errors.New("应用 ID 或密钥错误")
	}

	// 2. 检查应用状态
	if app.Status != 1 {
		return nil, errors.New("应用未启用")
	}

	// 3. 查询授权码
	var authCode models.OpenAuthCode
	err = l.svcCtx.DB.Where("code = ? AND app_id = ?", req.Code, req.AppID).First(&authCode).Error
	if err != nil {
		return nil, errors.New("授权码无效")
	}

	// 4. 检查是否已使用
	if authCode.Used {
		return nil, errors.New("授权码已使用")
	}

	// 5. 检查是否过期
	if authCode.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("授权码已过期")
	}

	// 6. 标记授权码为已使用
	authCode.Used = true
	l.svcCtx.DB.Save(&authCode)

	// 7. 生成 access_token 和 refresh_token
	accessToken := utils.GenerateAccessToken()
	refreshToken := utils.GenerateRefreshToken()

	// 8. 保存到数据库
	expiresIn := int64(7200) // 2小时
	accessTokenModel := models.OpenAccessToken{
		AppID:        req.AppID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    utils.GetTokenExpiry(expiresIn),
		Scope:        authCode.Scope,
	}
	l.svcCtx.DB.Create(&accessTokenModel)

	logx.Infof("OAuth 授权成功: app_id=%s, user_id=%s", req.AppID, authCode.UserID)

	return &types.GetTokenByCodeRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200,
		TokenType:    "Bearer",
	}, nil
}
