package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ResolveCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewResolveCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResolveCommentLogic {
	return &ResolveCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ResolveCommentLogic", ctx),
	}
}

func (l *ResolveCommentLogic) ResolveComment(in *doc_rpc.ResolveCommentReq) (*doc_rpc.ResolveCommentRes, error) {
	var row doc_models.DocComment
	if err := l.svcCtx.DB.Where("comment_id = ? AND deleted_at = ?",
		in.CommentId, zeroDeletedAt).First(&row).Error; err != nil {
		return nil, errNotFound
	}
	// 解决状态只挂在顶层评论上，回复随顶层一起关闭。
	if row.ParentID != "" {
		return nil, errInvalidArg
	}
	_, role, err := requireDoc(l.svcCtx.DB, row.DocID, in.OperatorUserId, roleViewer)
	if err != nil {
		return nil, err
	}
	actor := trim(in.OperatorUserId)
	if !canEdit(role) && trim(row.CreatedBy) != actor {
		return nil, errNoPermission
	}

	upd := map[string]any{"resolved": in.Resolved, "resolved_by": "", "resolved_at": ""}
	if in.Resolved {
		upd["resolved_by"] = actor
		upd["resolved_at"] = nowStr()
	}
	if err := l.svcCtx.DB.Model(&doc_models.DocComment{}).
		Where("comment_id = ?", row.CommentID).Updates(upd).Error; err != nil {
		return nil, err
	}
	auditJSON(l.svcCtx.DB, actor, "comment.resolve", "doc", row.DocID,
		map[string]any{"commentId": row.CommentID, "resolved": in.Resolved})
	return &doc_rpc.ResolveCommentRes{}, nil
}
