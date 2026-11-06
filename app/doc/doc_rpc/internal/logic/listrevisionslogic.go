package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ListRevisionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListRevisionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRevisionsLogic {
	return &ListRevisionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListRevisionsLogic", ctx),
	}
}

func (l *ListRevisionsLogic) ListRevisions(in *doc_rpc.ListRevisionsReq) (*doc_rpc.ListRevisionsRes, error) {
	db := l.svcCtx.DB.Model(&doc_models.DocRevision{}).Where("doc_id = ?", in.DocId)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "统计版本失败", Data: err})
		return nil, err
	}

	page, pageSize := normalizePage(in.Page, in.PageSize)
	var rows []doc_models.DocRevision
	if err := db.Order("version desc, id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "列出版本失败", Data: err})
		return nil, err
	}

	items := make([]*doc_rpc.RevisionItem, 0, len(rows))
	for i := range rows {
		items = append(items, toRevisionItem(&rows[i]))
	}
	return &doc_rpc.ListRevisionsRes{Total: total, Revisions: items}, nil
}
