package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ListTrashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListTrashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTrashLogic {
	return &ListTrashLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListTrashLogic", ctx),
	}
}

func (l *ListTrashLogic) ListTrash(in *doc_rpc.ListTrashReq) (*doc_rpc.ListTrashRes, error) {
	uid := trim(in.UserId)
	if uid == "" {
		return &doc_rpc.ListTrashRes{}, nil
	}
	q := l.svcCtx.DB.Model(&doc_models.DocTrash{}).Where("deleted_by = ?", uid)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}
	page, size := pageWindow(in.Page, in.PageSize)
	var rows []doc_models.DocTrash
	if err := q.Order("deleted_at desc, id desc").
		Offset((page - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]*doc_rpc.TrashItem, 0, len(rows))
	for i := range rows {
		items = append(items, toTrashItem(&rows[i]))
	}
	return &doc_rpc.ListTrashRes{Total: total, Items: items}, nil
}
