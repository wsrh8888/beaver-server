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
	// 1. 检查应用是否存在且属于当前用户
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	// 2. 查询或创建 OAuth 配置
	var oauthConfig open_models.OpenOAuthConfig
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&oauthConfig).Error; err != nil {
		// 不存在则创建
		oauthConfig = open_models.OpenOAuthConfig{
			AppID:           req.AppID,
			RedirectURIs:    "[]",
			Scopes:          "[]",
			CustomLogo:      "",
			CustomTitle:     "",
			CustomColor:     "",
			EnablePKCE:      0,
			TokenExpiration: 7200,
			Status:          1,
		}
	}

	// 3. 更新字段（只更新传入的字段）
	if req.RedirectURIs != nil {
		redirectURIsJSON, _ := json.Marshal(req.RedirectURIs)
		oauthConfig.RedirectURIs = string(redirectURIsJSON)
	}
	if req.Scopes != nil {
		scopesJSON, _ := json.Marshal(req.Scopes)
		oauthConfig.Scopes = string(scopesJSON)
	}
	if req.CustomLogo != "" {
		oauthConfig.CustomLogo = req.CustomLogo
	}
	if req.CustomTitle != "" {
		oauthConfig.CustomTitle = req.CustomTitle
	}
	if req.CustomColor != "" {
		oauthConfig.CustomColor = req.CustomColor
	}
	if req.EnablePKCE != nil {
		if *req.EnablePKCE {
			oauthConfig.EnablePKCE = 1
		} else {
			oauthConfig.EnablePKCE = 0
		}
	}
	if req.TokenExpiration != nil {
		oauthConfig.TokenExpiration = *req.TokenExpiration
	}
	if req.Status != nil {
		oauthConfig.Status = *req.Status
	}

	// 4. 保存配置
	if oauthConfig.ID == 0 {
		// 新建
		if err := l.svcCtx.DB.Create(&oauthConfig).Error; err != nil {
			return nil, errors.New("创建 OAuth 配置失败")
		}
	} else {
		// 更新
		if err := l.svcCtx.DB.Save(&oauthConfig).Error; err != nil {
			return nil, errors.New("更新 OAuth 配置失败")
		}
	}

	return &types.UpdateOAuthConfigRes{}, nil
}
