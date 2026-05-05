package app

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppScopesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取应用权限列表
func NewGetAppScopesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppScopesLogic {
	return &GetAppScopesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppScopesLogic) GetAppScopes(req *types.GetAppScopesReq) (resp *types.GetAppScopesRes, err error) {
	// todo: add your logic here and delete this line

	return
}
