package app

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOAuthConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新 OAuth 配置
func NewUpdateOAuthConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOAuthConfigLogic {
	return &UpdateOAuthConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOAuthConfigLogic) UpdateOAuthConfig(req *types.UpdateOAuthConfigReq) (resp *types.UpdateOAuthConfigRes, err error) {
	// todo: add your logic here and delete this line

	return
}
