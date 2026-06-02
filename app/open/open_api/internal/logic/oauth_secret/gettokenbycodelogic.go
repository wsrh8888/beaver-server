package oauth_secret

import (
	"context"

	"beaver/app/open/open_api/internal/logic/oauthutil"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenByCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTokenByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenByCodeLogic {
	return &GetTokenByCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenByCodeLogic) GetTokenByCode(req *types.GetTokenByCodeReq) (resp *types.GetTokenByCodeRes, err error) {
	if _, err := oauthutil.VerifyAppForCodeExchange(l.svcCtx.DB, req.AppID, req.AppSecret); err != nil {
		return nil, err
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

	return &types.GetTokenByCodeRes{
		AccessToken:  rpcResp.AccessToken,
		RefreshToken: rpcResp.RefreshToken,
		ExpiresIn:    rpcResp.ExpiresIn,
		TokenType:    rpcResp.TokenType,
		Scope:        scope,
	}, nil
}
