package logic

import (
	"context"
	"time"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DeleteCommentLogic", ctx),
	}
}

func (l *DeleteCommentLogic) DeleteComment(in *doc_rpc.DeleteCommentReq) (*doc_rpc.DeleteCommentRes, error) {
	var row doc_models.DocComment
	if err := l.svcCtx.DB.Where("comment_id = ? AND deleted_at = ?",
		in.CommentId, zeroDeletedAt).First(&row).Error; err != nil {
		return nil, errNotFound
	}
	_, role, err := requireDoc(l.svcCtx.DB, row.DocID, in.OperatorUserId, roleViewer)
	if err != nil {
		return nil, err
	}
	actor := trim(in.OperatorUserId)
	if !canManage(role) && trim(row.CreatedBy) != actor {
		return nil, errNoPermission
	}

	now := time.Now()
	if err := l.svcCtx.DB.Model(&doc_models.DocComment{}).
		Where("comment_id = ?", row.CommentID).Update("deleted_at", now).Error; err != nil {
		return nil, err
	}
	if row.ParentID == "" {
		_ = l.svcCtx.DB.Model(&doc_models.DocComment{}).
			Where("parent_id = ?", row.CommentID).Update("deleted_at", now).Error
	}
	auditJSON(l.svcCtx.DB, actor, "comment.delete", "doc", row.DocID,
		map[string]any{"commentId": row.CommentID})
	return &doc_rpc.DeleteCommentRes{}, nil
}
