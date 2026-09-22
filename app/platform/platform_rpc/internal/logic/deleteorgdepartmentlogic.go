package logic

import (
	"context"

	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type DeleteOrgDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteOrgDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOrgDepartmentLogic {
	return &DeleteOrgDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DeleteOrgDepartmentLogic", ctx),
	}
}

func (l *DeleteOrgDepartmentLogic) DeleteOrgDepartment(in *platform_rpc.DeleteOrgDepartmentReq) (*platform_rpc.DeleteOrgDepartmentRes, error) {
	// todo: add your logic here and delete this line

	return &platform_rpc.DeleteOrgDepartmentRes{}, nil
}
