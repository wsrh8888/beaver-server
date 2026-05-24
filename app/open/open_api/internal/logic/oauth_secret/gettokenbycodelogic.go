package oauth_secret

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenByCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用授权码换取 Token（支持 PKCE）
func NewGetTokenByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenByCodeLogic {
	return &GetTokenByCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenByCodeLogic) GetTokenByCode(req *types.GetTokenByCodeReq) (resp *types.GetTokenByCodeRes, err error) {
	// todo: add your logic here and delete this line

	return
}
