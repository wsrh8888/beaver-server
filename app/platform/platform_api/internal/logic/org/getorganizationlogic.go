package org

import (
	"context"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"

	beaverlog "beaver/utils/beaverlog"
)

// goctl 会把 logx 写进 imports，占位避免未使用导入。
var _ = logx.WithContext

type GetOrganizationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 查询当前组织
func NewGetOrganizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrganizationLogic {
	return &GetOrganizationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetOrganization", ctx),
	}
}

func (l *GetOrganizationLogic) GetOrganization(req *types.GetOrganizationReq) (resp *types.GetOrganizationRes, err error) {
	// todo: add your logic here and delete this line

	return
}
