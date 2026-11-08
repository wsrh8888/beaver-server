package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type CreateDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDocLogic {
	return &CreateDocLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateDocLogic) CreateDoc(req *types.CreateDocReq) (resp *types.CreateDocRes, err error) {
	res, err := l.svcCtx.Doc.CreateDoc(l.ctx, &doc_rpc.CreateDocReq{
		SpaceId:     req.SpaceId,
		NodeId:      req.NodeId,
		Title:       req.Title,
		ContentJson: req.ContentJson,
		OwnerUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.CreateDocRes{DocId: res.DocId}, nil
}
