package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListNodesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListNodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNodesLogic {
	return &ListNodesLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListNodesLogic) ListNodes(req *types.ListNodesReq) (resp *types.ListNodesRes, err error) {
	res, err := l.svcCtx.Doc.ListNodes(l.ctx, &doc_rpc.ListNodesReq{
		SpaceId: req.SpaceId,
		UserId:  req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.ListNodesRes{Nodes: convNodeItems(res.Nodes)}, nil
}
