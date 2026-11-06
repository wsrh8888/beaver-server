package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type GetDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDocLogic {
	return &GetDocLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetDocLogic) GetDoc(req *types.GetDocReq) (resp *types.GetDocRes, err error) {
	res, err := l.svcCtx.Doc.GetDoc(l.ctx, &doc_rpc.GetDocReq{
		DocId:  req.DocId,
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.GetDocRes{Doc: convDocItem(res.Doc)}, nil
}
