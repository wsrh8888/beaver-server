package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ListTemplatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTemplatesLogic {
	return &ListTemplatesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListTemplatesLogic", ctx),
	}
}

func (l *ListTemplatesLogic) ListTemplates(in *doc_rpc.ListTemplatesReq) (*doc_rpc.ListTemplatesRes, error) {
	uid := trim(in.UserId)
	q := l.svcCtx.DB.Model(&doc_models.DocTemplate{}).Where("status = ?", 1).
		Where("scope = ? OR scope = ? OR owner_user_id = ?", "builtin", "team", uid)
	if cat := trim(in.Category); cat != "" {
		q = q.Where("category = ?", cat)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}
	page, size := pageWindow(in.Page, in.PageSize)
	var rows []doc_models.DocTemplate
	if err := q.Order("use_count desc, id desc").
		Offset((page - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]*doc_rpc.TemplateItem, 0, len(rows))
	for i := range rows {
		items = append(items, toTemplateItem(&rows[i]))
	}
	return &doc_rpc.ListTemplatesRes{Total: total, Templates: items}, nil
}
