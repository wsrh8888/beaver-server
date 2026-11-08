package logic

import (
	"context"

	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type GetDocContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetDocContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDocContentLogic {
	return &GetDocContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetDocContentLogic", ctx),
	}
}

func (l *GetDocContentLogic) GetDocContent(in *doc_rpc.GetDocContentReq) (*doc_rpc.GetDocContentRes, error) {
	doc, err := loadDoc(l.svcCtx.DB, in.DocId)
	if err != nil {
		return nil, errNotFound
	}
	return &doc_rpc.GetDocContentRes{
		DocId:       doc.DocID,
		Title:       doc.Title,
		ContentText: doc.ContentText,
		Version:     int32(doc.Version),
		Status:      int32(doc.Status),
	}, nil
}
