package logic

import (
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/utils/jwts"
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (*types.RefreshTokenRes, error) {
	// 从Redis获取当前token
	tokenKey := fmt.Sprintf("login_%s", req.UserID)
	token, err := l.svcCtx.Redis.Get(tokenKey).Result()
	if err != nil {
		l.Logger.Errorf("获取用户token失败: %v", err)
		return nil, fmt.Errorf("token已失效")
	}

	// 验证token
	claims, err := jwts.ParseToken(token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		l.Logger.Errorf("解析token失败: %v", err)
		return nil, fmt.Errorf("token已失效")
	}

	// 生成新token
	newToken, err := jwts.GenToken(jwts.JwtPayLoad{
		Nickname: claims.Nickname,
		UserID:   claims.UserID,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		l.Logger.Errorf("生成新token失败: %v", err)
		return nil, fmt.Errorf("刷新token失败")
	}

	// 更新Redis中的token
	if err := l.svcCtx.Redis.Set(tokenKey, newToken, 24*time.Hour).Err(); err != nil {
		l.Logger.Errorf("更新token失败: %v", err)
		return nil, fmt.Errorf("刷新token失败")
	}

	// 记录刷新日志
	l.Logger.Infof("用户 %s 刷新token成功,时间: %s", req.UserID, time.Now().Format("2006-01-02 15:04:05"))

	return &types.RefreshTokenRes{
		Token: newToken,
	}, nil
}
