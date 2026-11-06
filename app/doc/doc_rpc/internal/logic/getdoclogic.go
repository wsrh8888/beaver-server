package logic

import (
	"context"

	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
)

type GetDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDocLogic {
	return &GetDocLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetDocLogic", ctx),
	}
}

func (l *GetDocLogic) GetDoc(in *doc_rpc.GetDocReq) (*doc_rpc.GetDocRes, error) {
	doc, err := loadDoc(l.svcCtx.DB, in.DocId)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveDocRole(l.svcCtx.DB, doc.SpaceID, doc.DocID, in.UserId)
	if err != nil {
		return nil, err
	}
	if role == "" {
		return nil, errNoPermission
	}
	return &doc_rpc.GetDocRes{Doc: toDocItem(doc)}, nil
}
