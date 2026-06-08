package oauth_secret

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type GetTokenByCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewGetTokenByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenByCodeLogic {
	return &GetTokenByCodeLogic{
		ctx:    ctx,
		logger: logger.New("get_token_by_code"),
		svcCtx: svcCtx,
	}
}

func (l *GetTokenByCodeLogic) GetTokenByCode(req *types.GetTokenByCodeReq) (resp *types.GetTokenByCodeRes, err error) {
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}
	if req.AppSecret == "" {
		return nil, errors.New("appSecret 不能为空")
	}
	if req.Code == "" {
		return nil, errors.New("code 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND app_secret = ? AND status = ?", req.AppID, req.AppSecret, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或凭证错误")
	}

	rpcResp, err := l.svcCtx.OpenRpc.ExchangeToken(l.ctx, &open_rpc.ExchangeTokenReq{
		AppId: req.AppID,
		Code:  req.Code,
	})
	if err != nil {
		logx.Errorf("ExchangeToken RPC 调用失败: %v", err)
		return nil, err
	}

	scope := ""
	var tokenRecord open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("token = ?", rpcResp.AccessToken).First(&tokenRecord).Error; err == nil {
		scope = tokenRecord.Scope
	}

	l.logger.Info(model.LogMsg{
		Text: "OAuth换发Token成功",
		Data: map[string]interface{}{
			"appId": req.AppID,
		},
	})

	return &types.GetTokenByCodeRes{
		AccessToken:  rpcResp.AccessToken,
		RefreshToken: rpcResp.RefreshToken,
		ExpiresIn:    rpcResp.ExpiresIn,
		TokenType:    rpcResp.TokenType,
		Scope:        scope,
	}, nil
}
