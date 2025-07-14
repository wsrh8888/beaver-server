package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/utils"
	"beaver/utils/device"
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

	// 从context获取User-Agent
	userAgent := l.ctx.Value("user-agent")
	logx.Infof("User-Agent: %s", userAgent)
	var deviceType string
	if userAgent == nil {
		deviceType = "unknown"
		logx.Errorf("User-Agent为空，用户: %s", claims.UserID)
	} else {
		deviceType = device.GetDeviceType(userAgent.(string))
		logx.Infof("认证设备类型识别 - 用户: %s, User-Agent: %s, 识别结果: %s", claims.UserID, userAgent.(string), deviceType)
	}

	// 直接构建特定设备类型的key
	key := fmt.Sprintf("login_%s_%s", claims.UserID, deviceType)
	loginInfoStr, err := l.svcCtx.Redis.Get(key).Result()
	if err != nil {
		logx.Errorf("获取登录信息失败: %v, Key: %s", err, key)
		return nil, errors.New("token已失效 " + key + " " + err.Error())
	}

	// 解析登录信息并验证token
	var loginInfo map[string]interface{}
	if err := json.Unmarshal([]byte(loginInfoStr), &loginInfo); err != nil {
		return nil, errors.New("登录信息格式错误")
	}

	storedToken, ok := loginInfo["token"].(string)
	if !ok || storedToken != req.Token {
		return nil, errors.New("token已失效或不匹配")
	}

	// 更新最后活跃时间
	loginInfo["last_active"] = time.Now().Format("2006-01-02 15:04:05")
	updatedInfo, _ := json.Marshal(loginInfo)
	l.svcCtx.Redis.Set(key, string(updatedInfo), time.Hour*48)

	resp = &types.AuthenticationRes{
		UserID: claims.UserID,
	}
	return resp, nil
}
