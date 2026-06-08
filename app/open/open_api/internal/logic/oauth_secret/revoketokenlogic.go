package oauth_secret

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
)


type RevokeTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewRevokeTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeTokenLogic {
	return &RevokeTokenLogic{
		ctx:    ctx,
		logger: logger.New("revoke_token"),
		svcCtx: svcCtx,
	}
}

func (l *RevokeTokenLogic) RevokeToken(req *types.RevokeTokenReq, appID, appSecret string) (resp *types.RevokeTokenRes, err error) {
	if req.Token == "" {
		return nil, errors.New("token 不能为空")
	}
	if appID == "" || appSecret == "" {
		return nil, errors.New("应用凭证不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND app_secret = ? AND status = ?", appID, appSecret, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或凭证错误")
	}

	var tokenRecord open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("token = ? OR refresh_token = ?", req.Token, req.Token).First(&tokenRecord).Error; err != nil {
		return nil, errors.New("令牌不存在")
	}
	if tokenRecord.AppID != appID {
		return nil, errors.New("无权撤销该令牌")
	}

	result := l.svcCtx.DB.Where("id = ?", tokenRecord.ID).Delete(&open_models.OpenOAuthToken{})
	if result.Error != nil {
		return nil, errors.New("撤销令牌失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "OAuth撤销Token成功",
		Data: map[string]interface{}{
			"appId": appID,
		},
	})

	return &types.RevokeTokenRes{Success: true}, nil
}
