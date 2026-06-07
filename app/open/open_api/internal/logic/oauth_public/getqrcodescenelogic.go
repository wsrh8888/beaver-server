package oauth_public

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	oauthmiddle "beaver/app/open/open_api/internal/middle/oauth"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/constants"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQrCodeSceneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetQrCodeSceneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQrCodeSceneLogic {
	return &GetQrCodeSceneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQrCodeSceneLogic) GetQrCodeScene(req *types.GetQrCodeSceneReq) (resp *types.GetQrCodeSceneRes, err error) {
	qrCode, err := l.svcCtx.OAuth.LoadScene(req.SceneID)
	if err != nil {
		return nil, err
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND status = ?", qrCode.AppID, 1).First(&app).Error; err != nil {
		return nil, err
	}

	expireIn := int64(time.Until(qrCode.ExpiresAt).Seconds())
	if expireIn < 0 {
		expireIn = 0
	}

	var oauthConfig open_models.OpenAppOAuth
	scopeStr := ""
	if err := l.svcCtx.DB.Where("app_id = ?", qrCode.AppID).First(&oauthConfig).Error; err == nil && oauthConfig.SupportedScopes != "" {
		scopeStr = oauthConfig.SupportedScopes
	} else {
		scopes := []string{
			string(constants.ScopeUserProfileRead),
			string(constants.ScopeUserAvatarRead),
		}
		data, _ := json.Marshal(scopes)
		scopeStr = string(data)
	}

	return &types.GetQrCodeSceneRes{
		SceneID:  qrCode.SceneID,
		AppID:    qrCode.AppID,
		AppName:  app.Name,
		AppIcon:  app.Icon,
		Status:   oauthmiddle.QrStatusText(qrCode.Status),
		ExpireIn: expireIn,
		Scopes:   parseScopeList(scopeStr),
	}, nil
}

func parseScopeList(scopeStr string) []string {
	scopeStr = strings.TrimSpace(scopeStr)
	if scopeStr == "" {
		return nil
	}
	if strings.HasPrefix(scopeStr, "[") {
		var scopes []string
		if err := json.Unmarshal([]byte(scopeStr), &scopes); err == nil {
			return scopes
		}
	}
	return strings.FieldsFunc(scopeStr, func(r rune) bool {
		return r == ' ' || r == ','
	})
}
