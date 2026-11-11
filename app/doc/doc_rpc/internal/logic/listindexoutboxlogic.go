package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ListIndexOutboxLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListIndexOutboxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListIndexOutboxLogic {
	return &ListIndexOutboxLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListIndexOutboxLogic", ctx),
	}
}

func (l *ListIndexOutboxLogic) ListIndexOutbox(in *doc_rpc.ListIndexOutboxReq) (*doc_rpc.ListIndexOutboxRes, error) {
	status := trim(in.Status)
	if status == "" {
		status = "pending"
	}
	limit := int(in.Limit)
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	var rows []doc_models.DocIndexOutbox
	if err := l.svcCtx.DB.Where("status = ?", status).
		Order("id asc").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]*doc_rpc.IndexOutboxItem, 0, len(rows))
	for i := range rows {
		items = append(items, &doc_rpc.IndexOutboxItem{
			DocId:     rows[i].DocID,
			Version:   int32(rows[i].Version),
			Digest:    rows[i].Digest,
			Status:    rows[i].Status,
			Attempts:  int32(rows[i].Attempts),
			CreatedAt: rows[i].CreatedStr,
		})
	}
	return &doc_rpc.ListIndexOutboxRes{Items: items}, nil
}
