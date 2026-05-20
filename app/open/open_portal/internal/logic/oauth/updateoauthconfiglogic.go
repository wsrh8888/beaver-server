package oauth

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOAuthConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新 OAuth 配置
func NewUpdateOAuthConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOAuthConfigLogic {
	return &UpdateOAuthConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOAuthConfigLogic) UpdateOAuthConfig(req *types.UpdateOAuthConfigReq) (resp *types.UpdateOAuthConfigRes, err error) {
	// 查询应用信息
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		return nil, err
	}

	// 验证 UserID 是否为应用所有者
	if app.OwnerUserID != req.UserID {
		return nil, errors.New("无权修改此应用")
	}

	// 解析现有的 OAuthConfig
	var oauthConfig open_models.OAuthClientConfig
	if app.OauthConfig != "" {
		if err := json.Unmarshal([]byte(app.OauthConfig), &oauthConfig); err != nil {
			l.Errorf("解析现有 OAuthConfig 失败: %v", err)
			// 如果解析失败，使用空配置
			oauthConfig = open_models.OAuthClientConfig{}
		}
	}

	// 根据类型更新对应配置
	switch req.OAuthType {
	case "h5":
		var h5Config open_models.H5OAuthConfig
		if err := json.Unmarshal([]byte(req.Config), &h5Config); err != nil {
			return nil, errors.New("H5 配置格式错误")
		}
		oauthConfig.H5 = h5Config

	case "desktop":
		var desktopConfig open_models.DesktopOAuthConfig
		if err := json.Unmarshal([]byte(req.Config), &desktopConfig); err != nil {
			return nil, errors.New("桌面端配置格式错误")
		}
		oauthConfig.Desktop = desktopConfig

	case "mobile":
		var mobileConfig open_models.MobileOAuthConfig
		if err := json.Unmarshal([]byte(req.Config), &mobileConfig); err != nil {
			return nil, errors.New("移动端配置格式错误")
		}
		oauthConfig.Mobile = mobileConfig

	default:
		return nil, errors.New("不支持的 OAuth 类型")
	}

	// 序列化回 JSON
	configJSON, err := json.Marshal(oauthConfig)
	if err != nil {
		return nil, errors.New("配置序列化失败")
	}

	// 更新数据库
	if err := l.svcCtx.DB.Model(&app).Update("oauth_config", string(configJSON)).Error; err != nil {
		return nil, err
	}

	return &types.UpdateOAuthConfigRes{}, nil
}
