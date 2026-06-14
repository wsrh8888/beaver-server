package auth_public

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/middleware/ua"
	"beaver/utils/device"
	"beaver/utils/jwts"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type OAuthCodeLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewOAuthCodeLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OAuthCodeLoginLogic {
	return &OAuthCodeLoginLogic{
		ctx:    ctx,
		logger: logger.New("oauth_code_login"),
		svcCtx: svcCtx,
	}
}

func (l *OAuthCodeLoginLogic) OAuthCodeLogin(req *types.OAuthCodeLoginReq) (*types.OAuthCodeLoginRes, error) {
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}
	if req.Code == "" {
		return nil, errors.New("授权码不能为空")
	}

	profile := ua.Profile(l.ctx)
	preciseType := ua.DeviceType(l.ctx)
	deviceGroup := ua.DeviceGroup(l.ctx)

	rpcResp, err := l.svcCtx.OpenRpc.ExchangeToken(l.ctx, &open_rpc.ExchangeTokenReq{
		AppId: req.AppID,
		Code:  req.Code,
	})
	if err != nil {
		logx.Errorf("OAuth code 换 token 失败: appId=%s, err=%v", req.AppID, err)
		return nil, errors.New("授权码无效或已过期")
	}
	if rpcResp.UserId == "" {
		return nil, errors.New("授权码换取用户信息失败")
	}

	userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: rpcResp.UserId,
	})
	if err != nil || userRes.UserInfo == nil {
		return nil, errors.New("用户不存在")
	}
	userInfo := userRes.UserInfo

	key := fmt.Sprintf("user_authentication_session:%s:%s", userInfo.UserId, deviceGroup)
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		NickName: userInfo.NickName,
		UserID:   userInfo.UserId,
		DeviceID: req.DeviceID,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return nil, errors.New("服务内部异常")
	}

	loginInfo, _ := json.Marshal(map[string]interface{}{
		"token": token, "device_id": req.DeviceID, "device_type": preciseType,
		"device_group": deviceGroup, "login_time": time.Now().Format("2006-01-02 15:04:05"),
	})
	if err = l.svcCtx.Redis.Set(key, string(loginInfo), time.Hour*time.Duration(l.svcCtx.Config.Auth.AccessExpire)).Err(); err != nil {
		return nil, errors.New("服务内部异常")
	}

	_ = device.UpsertOnLogin(l.svcCtx.DB, userInfo.UserId, req.DeviceID, profile, req.ClientIP)

	l.logger.Info(model.LogMsg{
		Text: "OAuth 授权码登录成功",
		Data: map[string]interface{}{
			"userId":      userInfo.UserId,
			"appId":       req.AppID,
			"deviceGroup": deviceGroup,
		},
	})

	return &types.OAuthCodeLoginRes{Token: token, UserID: userInfo.UserId}, nil
}
