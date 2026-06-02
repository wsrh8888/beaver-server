package oauth_secret

import (
	"context"

	"beaver/app/open/open_api/internal/logic/oauthutil"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRevokeTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeTokenLogic {
	return &RevokeTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeTokenLogic) RevokeToken(req *types.RevokeTokenReq) (resp *types.RevokeTokenRes, err error) {
	if err := oauthutil.RevokeOAuthToken(l.svcCtx.DB, req.Token); err != nil {
		return nil, err
	}

	return &types.RevokeTokenRes{Success: true}, nil
}
