package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("CreateCommentLogic", ctx),
	}
}

func (l *CreateCommentLogic) CreateComment(in *doc_rpc.CreateCommentReq) (*doc_rpc.CreateCommentRes, error) {
	doc, _, err := requireDoc(l.svcCtx.DB, in.DocId, in.OperatorUserId, roleCommenter)
	if err != nil {
		return nil, err
	}
	content := trim(in.Content)
	if content == "" {
		return nil, errInvalidArg
	}
	actor := trim(in.OperatorUserId)

	parentID := trim(in.ParentId)
	if parentID != "" {
		var parent doc_models.DocComment
		if err := l.svcCtx.DB.Where("comment_id = ? AND doc_id = ? AND deleted_at = ?",
			parentID, doc.DocID, zeroDeletedAt).First(&parent).Error; err != nil {
			return nil, errNotFound
		}
		// 回复统一挂到顶层评论下，保持两层结构。
		if parent.ParentID != "" {
			parentID = parent.ParentID
		}
	}

	row := doc_models.DocComment{
		CommentID: genID("cmt"),
		DocID:     doc.DocID,
		Epoch:     doc.Epoch,
		ParentID:  parentID,
		AnchorID:  trim(in.AnchorId),
		Quote:     firstRunes(in.Quote, 500),
		Content:   content,
		Mentions:  trim(in.Mentions),
		CreatedBy: actor,
		DeletedAt: sentinelDeletedAt,
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		return nil, err
	}

	for _, uid := range splitMentions(row.Mentions) {
		notify(l.svcCtx.DB, uid, doc.DocID, "mention", actor,
			"有人在《"+doc.Title+"》中提到了你", content)
	}
	if parentID != "" {
		var root doc_models.DocComment
		if err := l.svcCtx.DB.Where("comment_id = ?", parentID).First(&root).Error; err == nil {
			notify(l.svcCtx.DB, root.CreatedBy, doc.DocID, "comment_reply", actor,
				"《"+doc.Title+"》有新的回复", content)
		}
	} else {
		notify(l.svcCtx.DB, doc.OwnerUserID, doc.DocID, "comment", actor,
			"《"+doc.Title+"》有新评论", content)
	}
	auditJSON(l.svcCtx.DB, actor, "comment.create", "doc", doc.DocID,
		map[string]any{"commentId": row.CommentID, "parentId": parentID})

	return &doc_rpc.CreateCommentRes{Comment: toCommentItem(&row, nil)}, nil
}
