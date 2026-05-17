package oauth

import (
	"context"
	"encoding/json"
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

	resp = &types.GetOAuthConfigRes{
		OAuthType: req.OAuthType,
	}

	// 解析 OAuthConfig JSON
	if app.OauthConfig == "" {
		return resp, nil
	}

	var oauthConfig open_models.OAuthClientConfig
	if err := json.Unmarshal([]byte(app.OauthConfig), &oauthConfig); err != nil {
		l.Errorf("解析 OAuthConfig 失败: %v", err)
		return resp, nil
	}

	// 根据类型返回对应配置
	switch req.OAuthType {
	case "h5":
		resp.H5Config = &types.H5OAuthConfigInfo{
			Enabled:      oauthConfig.H5.Enabled,
			RedirectURIs: oauthConfig.H5.RedirectURIs,
			JsSdkDomains: oauthConfig.H5.JsSdkDomains,
		}
	case "desktop":
		// 生成授权页面 URL
		authPageURL := ""
		if oauthConfig.Desktop.Enabled && oauthConfig.Desktop.CustomScheme != "" {
			oauthBaseUrl := l.svcCtx.Config.OAuth.BaseUrl
			redirectURI := oauthConfig.Desktop.CustomScheme
			authPageURL = fmt.Sprintf("%s/desktop/auth?appId=%s&redirectUri=%s",
				oauthBaseUrl, req.AppID, url.QueryEscape(redirectURI))
		}

		resp.DesktopConfig = &types.DesktopOAuthConfigInfo{
			Enabled:      oauthConfig.Desktop.Enabled,
			CustomScheme: oauthConfig.Desktop.CustomScheme,
			AuthPageURL:  authPageURL,
		}
	case "mobile":
		resp.MobileConfig = &types.MobileOAuthConfigInfo{
			Enabled:            oauthConfig.Mobile.Enabled,
			IOSBundleID:        oauthConfig.Mobile.IOSBundleID,
			AndroidPackageName: oauthConfig.Mobile.AndroidPackageName,
			UniversalLink:      oauthConfig.Mobile.UniversalLink,
			CustomScheme:       oauthConfig.Mobile.CustomScheme,
		}
	}

	return resp, nil
}
