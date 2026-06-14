package auth_public

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/jwts"

	"github.com/zeromicro/go-zero/core/logx"
)

type OAuthLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOAuthLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OAuthLoginLogic {
	return &OAuthLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OAuthLoginLogic) OAuthLogin(req *types.OAuthLoginReq) (resp *types.OAuthLoginRes, err error) {
	if req.Code == "" {
		return nil, errors.New("授权码不能为空")
	}

	appID := l.svcCtx.Config.PortalOAuth.AppId
	if appID == "" {
		return nil, errors.New("门户 OAuth 未配置")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		logx.Errorf("查询门户 OAuth 应用失败: appId=%s, err=%v", appID, err)
		return nil, errors.New("门户 OAuth 应用不存在")
	}

	rpcResp, err := l.svcCtx.OpenRpc.ExchangeToken(l.ctx, &open_rpc.ExchangeTokenReq{
		AppId: appID,
		Code:  req.Code,
	})
	if err != nil {
		logx.Errorf("code 换 token 失败: appId=%s, err=%v", appID, err)
		return nil, errors.New("授权码换取令牌失败")
	}
	if rpcResp.AccessToken == "" {
		return nil, errors.New("授权码换取令牌失败")
	}

	var tokenRecord open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("token = ?", rpcResp.AccessToken).First(&tokenRecord).Error; err != nil {
		logx.Errorf("查询 OAuth token 失败: err=%v", err)
		return nil, errors.New("获取用户信息失败")
	}

	userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: tokenRecord.UserID,
	})
	if err != nil || userRes.UserInfo == nil {
		logx.Errorf("查询用户信息失败: userId=%s, err=%v", tokenRecord.UserID, err)
		return nil, errors.New("获取用户信息失败")
	}

	secretKey := l.svcCtx.Config.Auth.AccessSecret
	expireHours := l.svcCtx.Config.Auth.AccessExpire / 3600
	if expireHours == 0 {
		expireHours = 12
	}

	token, err := jwts.GenToken(jwts.JwtPayLoad{
		UserID:   userRes.UserInfo.UserId,
		NickName: userRes.UserInfo.NickName,
	}, secretKey, int(expireHours))
	if err != nil {
		logx.Errorf("生成 token 失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	expireAt := time.Now().Add(time.Duration(expireHours) * time.Hour).UnixMilli()
	logx.Infof("OAuth 登录成功: userId=%s, nickName=%s", userRes.UserInfo.UserId, userRes.UserInfo.NickName)

	return &types.OAuthLoginRes{
		Token:    token,
		UserID:   userRes.UserInfo.UserId,
		NickName: userRes.UserInfo.NickName,
		ExpireAt: expireAt,
	}, nil
}
