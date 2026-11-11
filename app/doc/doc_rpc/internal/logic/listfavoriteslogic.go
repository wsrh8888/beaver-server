package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ListFavoritesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFavoritesLogic {
	return &ListFavoritesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListFavoritesLogic", ctx),
	}
}

func (l *ListFavoritesLogic) ListFavorites(in *doc_rpc.ListFavoritesReq) (*doc_rpc.ListFavoritesRes, error) {
	uid := trim(in.UserId)
	if uid == "" {
		return &doc_rpc.ListFavoritesRes{}, nil
	}
	sub := l.svcCtx.DB.Model(&doc_models.DocFavorite{}).
		Select("doc_id").Where("user_id = ?", uid)
	q := visibleDocQuery(l.svcCtx.DB, uid).Where("doc_id IN (?)", sub)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}
	page, size := pageWindow(in.Page, in.PageSize)
	var rows []doc_models.DocDocument
	if err := q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, err
	}
	docs := make([]*doc_rpc.DocItem, 0, len(rows))
	for i := range rows {
		docs = append(docs, toDocItem(&rows[i]))
	}
	return &doc_rpc.ListFavoritesRes{Total: total, Docs: docs}, nil
}
