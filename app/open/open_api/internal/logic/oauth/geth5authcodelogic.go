package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"beaver/app/open/constants"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	util "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetH5AuthCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetH5AuthCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetH5AuthCodeLogic {
	return &GetH5AuthCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetH5AuthCodeLogic) GetH5AuthCode(req *types.GetH5AuthCodeReq) (resp *types.GetH5AuthCodeRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND status = ?", req.AppID, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或未启用")
	}

	var oauthConfig open_models.OpenAppOAuth
	scope := ""
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&oauthConfig).Error; err == nil && oauthConfig.SupportedScopes != "" {
		scope = oauthConfig.SupportedScopes
	} else {
		scopes := []string{
			string(constants.ScopeUserProfileRead),
			string(constants.ScopeUserAvatarRead),
		}
		data, _ := json.Marshal(scopes)
		scope = string(data)
	}

	const ttl = 5 * time.Minute
	authCode := util.NewV4().String()
	record := open_models.OpenOAuthCode{
		Code:      authCode,
		AppID:     req.AppID,
		UserID:    req.UserID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Scene:     "h5_sso",
	}
	if err := l.svcCtx.DB.Create(&record).Error; err != nil {
		logx.Errorf("生成 H5 authCode 失败: appId=%s, userId=%s, err=%v", req.AppID, req.UserID, err)
		return nil, errors.New("生成授权码失败")
	}

	logx.Infof("生成 H5 authCode 成功: appId=%s, userId=%s", req.AppID, req.UserID)

	return &types.GetH5AuthCodeRes{
		AuthCode: authCode,
		ExpireIn: int64(ttl.Seconds()),
	}, nil
}
