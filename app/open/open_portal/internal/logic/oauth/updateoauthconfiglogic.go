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

	// 查询或创建 OAuth 配置
	var oauth open_models.OpenAppOAuth
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&oauth).Error; err != nil {
		// 如果不存在，创建新记录
		oauth = open_models.OpenAppOAuth{
			AppID: req.AppID,
		}
	}

	// 根据类型更新对应配置
	switch req.OAuthType {
	case "h5":
		var h5Config open_models.H5OAuth
		if err := json.Unmarshal([]byte(req.Config), &h5Config); err != nil {
			return nil, errors.New("H5 配置格式错误")
		}
		oauth.H5 = &h5Config

	case "desktop":
		var desktopConfig open_models.DesktopOAuth
		if err := json.Unmarshal([]byte(req.Config), &desktopConfig); err != nil {
			return nil, errors.New("桌面端配置格式错误")
		}
		oauth.Desktop = &desktopConfig

	case "mobile":
		var mobileConfig open_models.MobileOAuth
		if err := json.Unmarshal([]byte(req.Config), &mobileConfig); err != nil {
			return nil, errors.New("移动端配置格式错误")
		}
		oauth.Mobile = &mobileConfig

	default:
		return nil, errors.New("不支持的 OAuth 类型")
	}

	// 保存或更新
	if oauth.ID == 0 {
		// 新建
		if err := l.svcCtx.DB.Create(&oauth).Error; err != nil {
			return nil, errors.New("保存 OAuth 配置失败")
		}
	} else {
		// 更新
		if err := l.svcCtx.DB.Save(&oauth).Error; err != nil {
			return nil, errors.New("更新 OAuth 配置失败")
		}
	}

	return &types.UpdateOAuthConfigRes{}, nil
}
