package oauth

import (
	"context"
	"encoding/json"

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
	// 1. 查询 OAuth 配置
	var oauthConfig open_models.OpenOAuthConfig
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&oauthConfig).Error; err != nil {
		// 如果不存在，返回默认配置
		return &types.GetOAuthConfigRes{
			Config: types.OAuthConfigInfo{
				AppID:           req.AppID,
				RedirectURIs:    []string{},
				Scopes:          []string{},
				CustomLogo:      "",
				CustomTitle:     "",
				CustomColor:     "",
				EnablePKCE:      false,
				TokenExpiration: 7200,
				Status:          1,
			},
		}, nil
	}

	// 2. 解析 JSON 字段
	var redirectURIs []string
	var scopes []string
	if oauthConfig.RedirectURIs != "" {
		json.Unmarshal([]byte(oauthConfig.RedirectURIs), &redirectURIs)
	}
	if oauthConfig.Scopes != "" {
		json.Unmarshal([]byte(oauthConfig.Scopes), &scopes)
	}

	// 3. 返回配置
	return &types.GetOAuthConfigRes{
		Config: types.OAuthConfigInfo{
			AppID:           oauthConfig.AppID,
			RedirectURIs:    redirectURIs,
			Scopes:          scopes,
			CustomLogo:      oauthConfig.CustomLogo,
			CustomTitle:     oauthConfig.CustomTitle,
			CustomColor:     oauthConfig.CustomColor,
			EnablePKCE:      oauthConfig.EnablePKCE == 1,
			TokenExpiration: oauthConfig.TokenExpiration,
			Status:          oauthConfig.Status,
		},
	}, nil
}
