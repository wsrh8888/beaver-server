package oauth

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenByCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用授权码换取 Token
func NewGetTokenByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenByCodeLogic {
	return &GetTokenByCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenByCodeLogic) GetTokenByCode(req *types.GetTokenByCodeReq) (resp *types.GetTokenByCodeRes, err error) {
	// 调用 OAuth RPC 服务
	rpcResp, err := l.svcCtx.OpenRpc.ExchangeToken(l.ctx, &open_rpc.ExchangeTokenReq{
		AppId: req.AppID,
		Code:  req.Code,
	})
	if err != nil {
		logx.Errorf("ExchangeToken RPC 调用失败: %v", err)
		return nil, err
	}

	return &types.GetTokenByCodeRes{
		AccessToken:  rpcResp.AccessToken,
		RefreshToken: rpcResp.RefreshToken,
		ExpiresIn:    rpcResp.ExpiresIn,
		TokenType:    rpcResp.TokenType,
	}, nil
}
