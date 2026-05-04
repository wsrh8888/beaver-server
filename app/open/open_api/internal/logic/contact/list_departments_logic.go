// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package contact

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDepartmentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取部门列表
func NewListDepartmentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDepartmentsLogic {
	return &ListDepartmentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDepartmentsLogic) ListDepartments(req *types.ListDepartmentsReq) (resp *types.ListDepartmentsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
