package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/utils"
	"beaver/utils/jwts"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthenticationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthenticationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationLogic {
	return &AuthenticationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthenticationLogic) Authentication(req *types.AuthenticationReq) (resp *types.AuthenticationRes, err error) {
	if utils.InListByRegex(l.svcCtx.Config.WhiteList, req.ValidPath) {
		logx.Infof("白名单请求：%s, %s", req.ValidPath, req.Token)
		return
	}
	if req.Token == "" {
		return nil, errors.New("token不能为空")
	}
	claims, err := jwts.ParseToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		logx.Errorf("解析token失败: %v", err)
		return nil, errors.New("认证失败")
	}
	key := fmt.Sprintf("login_%s", claims.UserID)
	token, err := l.svcCtx.Redis.Get(key).Result()
	if err != nil {
		logx.Errorf("获取Redis token失败: %v", err)
		return nil, errors.New("token已失效")
	}
	if token != req.Token {
		logx.Errorf("token不一致: Redis=%s, Request=%s", token, req.Token)
		return nil, errors.New("token已失效")
	}
	if err := l.svcCtx.Redis.Expire(key, time.Duration(l.svcCtx.Config.Auth.AccessExpire)*time.Second).Err(); err != nil {
		logx.Errorf("更新token过期时间失败: %v", err)
	}
	deviceKey := fmt.Sprintf("device_%s", claims.UserID)
	deviceInfo := map[string]interface{}{
		"last_active": time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := l.svcCtx.Redis.HMSet(deviceKey, deviceInfo).Err(); err != nil {
		logx.Errorf("更新设备信息失败: %v", err)
	}
	resp = &types.AuthenticationRes{
		UserID: claims.UserID,
	}
	fmt.Println(resp, "数据")
	return resp, nil
}
