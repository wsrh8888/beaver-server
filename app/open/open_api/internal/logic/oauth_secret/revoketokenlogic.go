package oauth_secret

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 撤销 Token（登出 / 解除授权）
func NewRevokeTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeTokenLogic {
	return &RevokeTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeTokenLogic) RevokeToken(req *types.RevokeTokenReq) (resp *types.RevokeTokenRes, err error) {
	// todo: add your logic here and delete this line

	return
}
