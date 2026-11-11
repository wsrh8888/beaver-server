package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ListRecentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListRecentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRecentLogic {
	return &ListRecentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListRecentLogic", ctx),
	}
}

func (l *ListRecentLogic) ListRecent(in *doc_rpc.ListRecentReq) (*doc_rpc.ListRecentRes, error) {
	uid := trim(in.UserId)
	if uid == "" {
		return &doc_rpc.ListRecentRes{}, nil
	}
	var visits []doc_models.DocRecentVisit
	if err := l.svcCtx.DB.Where("user_id = ?", uid).
		Order("visited_at desc, id desc").Find(&visits).Error; err != nil {
		return nil, err
	}
	if len(visits) == 0 {
		return &doc_rpc.ListRecentRes{}, nil
	}
	ids := make([]string, 0, len(visits))
	for i := range visits {
		ids = append(ids, visits[i].DocID)
	}
	var rows []doc_models.DocDocument
	if err := visibleDocQuery(l.svcCtx.DB, uid).Where("doc_id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	byID := make(map[string]*doc_models.DocDocument, len(rows))
	for i := range rows {
		byID[rows[i].DocID] = &rows[i]
	}
	page, size := pageWindow(in.Page, in.PageSize)
	start, end := (page-1)*size, page*size
	docs := make([]*doc_rpc.DocItem, 0, size)
	total := int64(0)
	for _, id := range ids {
		doc, ok := byID[id]
		if !ok {
			continue
		}
		if int(total) >= start && int(total) < end {
			docs = append(docs, toDocItem(doc))
		}
		total++
	}
	return &doc_rpc.ListRecentRes{Total: total, Docs: docs}, nil
}
