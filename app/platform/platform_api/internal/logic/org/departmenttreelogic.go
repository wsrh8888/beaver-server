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

type DepartmentTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 部门树，节点带直属成员人数
func NewDepartmentTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentTreeLogic {
	return &DepartmentTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DepartmentTree", ctx),
	}
}

func (l *DepartmentTreeLogic) DepartmentTree(req *types.DepartmentTreeReq) (resp *types.DepartmentTreeRes, err error) {
	// todo: add your logic here and delete this line

	return
}
