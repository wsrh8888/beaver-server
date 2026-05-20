package oauth

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOIDCUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户信息（OIDC 标准 userinfo 端点）
func NewGetOIDCUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOIDCUserInfoLogic {
	return &GetOIDCUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOIDCUserInfoLogic) GetOIDCUserInfo(req *types.GetOIDCUserInfoReq) (resp *types.GetOIDCUserInfoRes, err error) {
	// todo: add your logic here and delete this line

	return
}
