package security

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSecurityConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取安全配置
func NewGetSecurityConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSecurityConfigLogic {
	return &GetSecurityConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSecurityConfigLogic) GetSecurityConfig(req *types.GetSecurityConfigReq) (resp *types.GetSecurityConfigRes, err error) {
	// 查询应用信息
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		return nil, err
	}

	// 验证 UserID 是否为应用所有者
	if app.OwnerUserID != req.UserID {
		return nil, errors.New("无权查看此应用")
	}

	// 解析 IP 白名单
	var ipWhitelist []string
	if app.IPWhitelist != "" {
		json.Unmarshal([]byte(app.IPWhitelist), &ipWhitelist)
	}

	// 解析 H5 可信域名
	var trustedDomains []string
	if app.TrustedDomains != "" {
		json.Unmarshal([]byte(app.TrustedDomains), &trustedDomains)
	}

	// 从 OAuth 配置表中读取 redirectUris
	var oauthConfig open_models.OpenOAuthConfig
	var redirectUris []string
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&oauthConfig).Error; err == nil {
		if oauthConfig.RedirectURIs != "" {
			json.Unmarshal([]byte(oauthConfig.RedirectURIs), &redirectUris)
		}
	}

	return &types.GetSecurityConfigRes{
		Config: types.SecurityConfigInfo{
			AppID:          app.AppID,
			IPWhitelist:    ipWhitelist,
			RedirectURIs:   redirectUris,
			TrustedDomains: trustedDomains,
		},
	}, nil
}
