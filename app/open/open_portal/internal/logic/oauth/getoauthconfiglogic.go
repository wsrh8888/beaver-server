package oauth

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOAuthConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 OAuth 配置
func NewGetOAuthConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOAuthConfigLogic {
	return &GetOAuthConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOAuthConfigLogic) GetOAuthConfig(req *types.GetOAuthConfigReq) (resp *types.GetOAuthConfigRes, err error) {

	// 查询应用信息
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		return nil, err
	}

	// 验证 UserID 是否为应用所有者
	if app.OwnerUserID != req.UserID {
		return nil, errors.New("无权查看此应用")
	}

	resp = &types.GetOAuthConfigRes{}

	// 查询 OAuth 配置
	var oauth open_models.OpenAppOAuth
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&oauth).Error; err != nil {
		// 如果没有配置，返回空配置
		return resp, nil
	}

	// 返回 H5 配置
	if oauth.H5 != nil {
		resp.H5Config = &types.H5OAuthConfigInfo{
			Enabled:      oauth.H5.Enabled,
			RedirectURIs: oauth.H5.RedirectURIs,
			JsSdkDomains: oauth.H5.JsSdkDomains,
		}
	}

	// 返回桌面端配置
	if oauth.Desktop != nil && oauth.Desktop.Enabled {
		authPageURL := ""
		if oauth.Desktop.CustomScheme != "" {
			oauthBaseUrl := l.svcCtx.Config.OAuth.BaseUrl
			redirectURI := oauth.Desktop.CustomScheme
			authPageURL = fmt.Sprintf("%s/desktop/auth?appId=%s&redirectUri=%s",
				oauthBaseUrl, req.AppID, url.QueryEscape(redirectURI))
		}

		resp.DesktopConfig = &types.DesktopOAuthConfigInfo{
			Enabled:      oauth.Desktop.Enabled,
			CustomScheme: oauth.Desktop.CustomScheme,
			AuthPageURL:  authPageURL,
		}
	}

	// 返回移动端配置
	if oauth.Mobile != nil {
		resp.MobileConfig = &types.MobileOAuthConfigInfo{
			Enabled:            oauth.Mobile.Enabled,
			IOSBundleID:        oauth.Mobile.IOSBundleID,
			AndroidPackageName: oauth.Mobile.AndroidPackageName,
			UniversalLink:      oauth.Mobile.UniversalLink,
			CustomScheme:       oauth.Mobile.CustomScheme,
		}
	}

	return resp, nil
}
