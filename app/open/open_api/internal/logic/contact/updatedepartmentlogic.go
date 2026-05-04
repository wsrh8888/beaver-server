package contact

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新部门
func NewUpdateDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDepartmentLogic {
	return &UpdateDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDepartmentLogic) UpdateDepartment(req *types.UpdateDepartmentReq) (resp *types.UpdateDepartmentRes, err error) {
	// TODO: 部门功能暂未实现
	return nil, errors.New("部门功能暂未实现")
}
