package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ListDocsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListDocsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDocsLogic {
	return &ListDocsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListDocsLogic", ctx),
	}
}

func (l *ListDocsLogic) ListDocs(in *doc_rpc.ListDocsReq) (*doc_rpc.ListDocsRes, error) {
	q := visibleDocQuery(l.svcCtx.DB, in.UserId)
	if sid := trim(in.SpaceId); sid != "" {
		q = q.Where("space_id = ?", sid)
	}
	if dt := trim(in.DocType); dt != "" {
		q = q.Where("doc_type = ?", dt)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}
	page, size := pageWindow(in.Page, in.PageSize)
	var rows []doc_models.DocDocument
	if err := q.Order("updated_at desc, id desc").
		Offset((page - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, err
	}
	docs := make([]*doc_rpc.DocItem, 0, len(rows))
	for i := range rows {
		docs = append(docs, toDocItem(&rows[i]))
	}
	return &doc_rpc.ListDocsRes{Total: total, Docs: docs}, nil
}
