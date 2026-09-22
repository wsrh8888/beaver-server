package logic

import (
	"context"

	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type AddOrgMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewAddOrgMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddOrgMemberLogic {
	return &AddOrgMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("AddOrgMemberLogic", ctx),
	}
}

func (l *AddOrgMemberLogic) AddOrgMember(in *platform_rpc.AddOrgMemberReq) (*platform_rpc.AddOrgMemberRes, error) {
	// todo: add your logic here and delete this line

	return &platform_rpc.AddOrgMemberRes{}, nil
}
