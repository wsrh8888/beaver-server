package logic

import (
	"context"

	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type UpdateOrgDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateOrgDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrgDepartmentLogic {
	return &UpdateOrgDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("UpdateOrgDepartmentLogic", ctx),
	}
}

func (l *UpdateOrgDepartmentLogic) UpdateOrgDepartment(in *platform_rpc.UpdateOrgDepartmentReq) (*platform_rpc.UpdateOrgDepartmentRes, error) {
	// todo: add your logic here and delete this line

	return &platform_rpc.UpdateOrgDepartmentRes{}, nil
}
