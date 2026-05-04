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

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 刷新访问令牌
func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenRes, err error) {
	// 1. 查询旧的 Access Token
	var oldToken open_models.OpenAccessToken
	if err := l.svcCtx.DB.Where("refresh_token = ?", req.RefreshToken).First(&oldToken).Error; err != nil {
		return nil, errors.New("刷新令牌无效")
	}

	// 2. 检查是否过期
	if time.Now().Unix() > oldToken.ExpiresAt {
		return nil, errors.New("令牌已过期，请重新授权")
	}

	// 3. 生成新的 Access Token
	accessTokenBytes := make([]byte, 32)
	_, _ = rand.Read(accessTokenBytes)
	newAccessToken := hex.EncodeToString(accessTokenBytes)

	// 4. 生成新的 Refresh Token
	refreshTokenBytes := make([]byte, 32)
	_, _ = rand.Read(refreshTokenBytes)
	newRefreshToken := hex.EncodeToString(refreshTokenBytes)

	// 5. 删除旧令牌记录
	if err := l.svcCtx.DB.Delete(&oldToken).Error; err != nil {
		logx.Errorf("删除旧令牌失败: %v", err)
	}

	// 6. 保存新令牌
	now := time.Now()
	expiresAt := now.Add(2 * time.Hour).Unix()
	newTokenRecord := open_models.OpenAccessToken{
		AppID:        oldToken.AppID,
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
		Scope:        oldToken.Scope,
		UserID:       oldToken.UserID,
	}

	if err := l.svcCtx.DB.Create(&newTokenRecord).Error; err != nil {
		logx.Errorf("创建新令牌失败: %v", err)
		return nil, errors.New("刷新令牌失败")
	}

	return &types.RefreshTokenRes{
		AccessToken: newAccessToken,
		ExpiresIn:   7200,
	}, nil
}
