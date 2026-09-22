package logic

import (
	"context"

	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type RemoveOrgMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRemoveOrgMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveOrgMemberLogic {
	return &RemoveOrgMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("RemoveOrgMemberLogic", ctx),
	}
}

func (l *RemoveOrgMemberLogic) RemoveOrgMember(in *platform_rpc.RemoveOrgMemberReq) (*platform_rpc.RemoveOrgMemberRes, error) {
	// todo: add your logic here and delete this line

	return &platform_rpc.RemoveOrgMemberRes{}, nil
}
