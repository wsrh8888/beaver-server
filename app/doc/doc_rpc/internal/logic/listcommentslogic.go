package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ListCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListCommentsLogic", ctx),
	}
}

func (l *ListCommentsLogic) ListComments(in *doc_rpc.ListCommentsReq) (*doc_rpc.ListCommentsRes, error) {
	doc, _, err := requireDoc(l.svcCtx.DB, in.DocId, in.OperatorUserId, roleViewer)
	if err != nil {
		return nil, err
	}
	q := l.svcCtx.DB.Model(&doc_models.DocComment{}).
		Where("doc_id = ? AND deleted_at = ?", doc.DocID, zeroDeletedAt)
	if !in.IncludeResolved {
		q = q.Where("resolved = ?", false)
	}
	var rows []doc_models.DocComment
	if err := q.Order("id asc").Find(&rows).Error; err != nil {
		return nil, err
	}

	byID := make(map[string]*doc_rpc.CommentItem, len(rows))
	items := make([]*doc_rpc.CommentItem, 0, len(rows))
	for i := range rows {
		it := toCommentItem(&rows[i], nil)
		byID[rows[i].CommentID] = it
		if rows[i].ParentID == "" {
			items = append(items, it)
		}
	}
	for i := range rows {
		if rows[i].ParentID == "" {
			continue
		}
		if p := byID[rows[i].ParentID]; p != nil {
			if reply := byID[rows[i].CommentID]; reply != nil {
				p.Replies = append(p.Replies, reply)
			}
		}
	}
	return &doc_rpc.ListCommentsRes{Total: int64(len(items)), Comments: items}, nil
}
