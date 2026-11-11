package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type SearchDocsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSearchDocsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchDocsLogic {
	return &SearchDocsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("SearchDocsLogic", ctx),
	}
}

func (l *SearchDocsLogic) SearchDocs(in *doc_rpc.SearchDocsReq) (*doc_rpc.SearchDocsRes, error) {
	q := visibleDocQuery(l.svcCtx.DB, in.UserId)
	if kw := trim(in.Query); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("title LIKE ? OR content_text LIKE ?", like, like)
	}
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
	items := make([]*doc_rpc.DocSearchItem, 0, len(rows))
	for i := range rows {
		items = append(items, &doc_rpc.DocSearchItem{
			Doc:     toDocItem(&rows[i]),
			Snippet: snippetAround(rows[i].ContentText, in.Query, 160),
		})
	}
	return &doc_rpc.SearchDocsRes{Total: total, Items: items}, nil
}
