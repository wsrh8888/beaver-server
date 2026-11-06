package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type DeleteSpaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSpaceLogic {
	return &DeleteSpaceLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteSpaceLogic) DeleteSpace(req *types.DeleteSpaceReq) (resp *types.DeleteSpaceRes, err error) {
	_, err = l.svcCtx.Doc.DeleteSpace(l.ctx, &doc_rpc.DeleteSpaceReq{
		SpaceId:        req.SpaceId,
		OperatorUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.DeleteSpaceRes{}, nil
}
