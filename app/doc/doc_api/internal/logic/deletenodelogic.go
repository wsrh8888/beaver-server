package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type DeleteNodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNodeLogic {
	return &DeleteNodeLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteNodeLogic) DeleteNode(req *types.DeleteNodeReq) (resp *types.DeleteNodeRes, err error) {
	_, err = l.svcCtx.Doc.DeleteNode(l.ctx, &doc_rpc.DeleteNodeReq{
		NodeId:         req.NodeId,
		OperatorUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.DeleteNodeRes{}, nil
}
