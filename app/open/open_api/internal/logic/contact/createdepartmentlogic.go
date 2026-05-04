package contact

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建部门
func NewCreateDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDepartmentLogic {
	return &CreateDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDepartmentLogic) CreateDepartment(req *types.CreateDepartmentReq) (resp *types.CreateDepartmentRes, err error) {
	// TODO: 目前返回错误，需要根据实际的部门模型实现
	return nil, errors.New("部门功能暂未实现")
}
