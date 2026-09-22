package logic

import (
	"context"

	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type CreateOrgDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateOrgDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrgDepartmentLogic {
	return &CreateOrgDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("CreateOrgDepartmentLogic", ctx),
	}
}

func (l *CreateOrgDepartmentLogic) CreateOrgDepartment(in *platform_rpc.CreateOrgDepartmentReq) (*platform_rpc.CreateOrgDepartmentRes, error) {
	// todo: add your logic here and delete this line

	return &platform_rpc.CreateOrgDepartmentRes{}, nil
}
