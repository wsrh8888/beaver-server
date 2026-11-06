package logic

import (
	"context"

	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
)

type GetSpaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSpaceLogic {
	return &GetSpaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetSpaceLogic", ctx),
	}
}

func (l *GetSpaceLogic) GetSpace(in *doc_rpc.GetSpaceReq) (*doc_rpc.GetSpaceRes, error) {
	space, err := loadSpace(l.svcCtx.DB, in.SpaceId)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveSpaceRole(l.svcCtx.DB, space, in.UserId)
	if err != nil {
		return nil, err
	}
	if role == "" {
		return nil, errNoPermission
	}
	return &doc_rpc.GetSpaceRes{Space: toSpaceItem(space)}, nil
}
