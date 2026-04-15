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

	// 2. 检查应用是否过期
	if app.ExpireAt > 0 && app.ExpireAt < time.Now().UnixMilli() {
		return nil, errors.New("应用已过期")
	}

	// 3. 生成 access_token
	accessToken := uuid.New().String()
	refreshToken := uuid.New().String()
	now := time.Now().UnixMilli()
	expiresIn := int64(7200) // 2 小时
	expiresAt := now + expiresIn*1000

	// 4. 保存 token
	tokenRecord := open_models.OpenAccessToken{
		Model: open_models.Model{
			ID:        uuid.New().String(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		AppID:       req.AppID,
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
		Status:      1, // 有效
	}

	err = l.svcCtx.DB.Create(&tokenRecord).Error
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 5. 保存 refresh_token
	refreshRecord := open_models.OpenRefreshToken{
		Model: open_models.Model{
			ID:        uuid.New().String(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		AppID:        req.AppID,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		ExpiresAt:    now + 30*24*3600*1000, // 30 天
		Status:       1,
	}

	err = l.svcCtx.DB.Create(&refreshRecord).Error
	if err != nil {
		return nil, errors.New("生成刷新令牌失败")
	}

	// 6. 记录 API 调用日志
	apiLog := open_models.OpenAPILog{
		Model: open_models.Model{
			ID:        uuid.New().String(),
			CreatedAt: now,
		},
		AppID:      req.AppID,
		APIPath:    "/api/open/v1/auth/token",
		Method:     "POST",
		StatusCode: 200,
		RequestIP:  "",
	}
	l.svcCtx.DB.Create(&apiLog)

	return &types.GetTokenRes{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		RefreshToken: refreshToken,
	}, nil
}
