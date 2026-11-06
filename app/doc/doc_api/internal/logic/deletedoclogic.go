package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type DeleteDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDocLogic {
	return &DeleteDocLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteDocLogic) DeleteDoc(req *types.DeleteDocReq) (resp *types.DeleteDocRes, err error) {
	_, err = l.svcCtx.Doc.DeleteDoc(l.ctx, &doc_rpc.DeleteDocReq{
		DocId:          req.DocId,
		OperatorUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.DeleteDocRes{}, nil
}
