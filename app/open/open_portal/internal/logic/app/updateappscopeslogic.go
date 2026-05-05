package app

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppScopesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新应用权限
func NewUpdateAppScopesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppScopesLogic {
	return &UpdateAppScopesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAppScopesLogic) UpdateAppScopes(req *types.UpdateAppScopesReq) (resp *types.UpdateAppScopesRes, err error) {
	// todo: add your logic here and delete this line

	return
}
