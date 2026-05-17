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

type UpdateSecurityConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新安全配置
func NewUpdateSecurityConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSecurityConfigLogic {
	return &UpdateSecurityConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSecurityConfigLogic) UpdateSecurityConfig(req *types.UpdateSecurityConfigReq) (resp *types.UpdateSecurityConfigRes, err error) {
	// 验证应用是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		return nil, err
	}

	// 验证 UserID 是否为应用所有者
	if app.OwnerUserID != req.UserID {
		return nil, errors.New("无权操作此应用")
	}

	// 将 IP 白名单转换为 JSON 字符串
	ipWhitelistJSON := ""
	if len(req.IPWhitelist) > 0 {
		data, _ := json.Marshal(req.IPWhitelist)
		ipWhitelistJSON = string(data)
	}

	// 将 H5 可信域名转换为 JSON 字符串
	trustedDomainsJSON := ""
	if len(req.TrustedDomains) > 0 {
		data, _ := json.Marshal(req.TrustedDomains)
		trustedDomainsJSON = string(data)
	}

	// 更新应用表（只更新 IP 白名单和 H5 可信域名）
	appUpdates := map[string]interface{}{
		"ip_whitelist":    ipWhitelistJSON,
		"trusted_domains": trustedDomainsJSON,
	}

	if err := l.svcCtx.DB.Model(&open_models.OpenApp{}).Where("app_id = ?", req.AppID).Updates(appUpdates).Error; err != nil {
		return nil, err
	}

	return &types.UpdateSecurityConfigRes{}, nil
}
