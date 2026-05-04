package oauth

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuthorizeCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取授权码
func NewGetAuthorizeCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuthorizeCodeLogic {
	return &GetAuthorizeCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAuthorizeCodeLogic) GetAuthorizeCode(req *types.GetAuthorizeCodeReq) (resp *types.AuthorizeCodeRes, err error) {
	// todo: add your logic here and delete this line

	return
}
