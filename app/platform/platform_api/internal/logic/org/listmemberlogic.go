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

type ListMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 成员列表，可按部门过滤
func NewListMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMemberLogic {
	return &ListMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListMember", ctx),
	}
}

func (l *ListMemberLogic) ListMember(req *types.ListMemberReq) (resp *types.ListMemberRes, err error) {
	// todo: add your logic here and delete this line

	return
}
