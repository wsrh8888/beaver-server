package auth

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/google/uuid"
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
	// 1. 验证 app_id 和 app_secret
	var app open_models.OpenApp
	err = l.svcCtx.DB.Where("app_id = ? AND app_secret = ? AND status = ?", req.AppID, req.AppSecret, 1).First(&app).Error
	if err != nil {
		return nil, errors.New("应用 ID 或密钥错误")
	}

	// 2. 生成 access_token
	accessToken := uuid.New().String()
	refreshToken := uuid.New().String()
	now := time.Now()
	expiresIn := int64(7200) // 2 小时
	expiresAt := now.Add(time.Duration(expiresIn) * time.Second).UnixMilli()

	// 3. 保存 token
	tokenRecord := open_models.OpenAccessToken{
		AppID:        req.AppID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Scope:        "",
	}

	err = l.svcCtx.DB.Create(&tokenRecord).Error
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	return &types.GetTokenRes{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		RefreshToken: refreshToken,
	}, nil
}
