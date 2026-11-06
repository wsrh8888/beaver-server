package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
)

type RollbackRevisionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRollbackRevisionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RollbackRevisionLogic {
	return &RollbackRevisionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("RollbackRevisionLogic", ctx),
	}
}

func (l *RollbackRevisionLogic) RollbackRevision(in *doc_rpc.RollbackRevisionReq) (*doc_rpc.RollbackRevisionRes, error) {
	doc, err := loadDoc(l.svcCtx.DB, in.DocId)
	if err != nil {
		return nil, errNotFound
	}
	space, err := loadSpace(l.svcCtx.DB, doc.SpaceID)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveDocRole(l.svcCtx.DB, doc.SpaceID, doc.DocID, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canEdit(role) {
		return nil, errNoPermission
	}

	var rev doc_models.DocRevision
	if err := l.svcCtx.DB.Where("revision_id = ? AND doc_id = ?", in.RevisionId, doc.DocID).First(&rev).Error; err != nil {
		return nil, errNotFound
	}

	// 回滚 = 用该版本正文生成新版本（保留历史线性，不删中间版本）。
	contentText := extractTextFromJSON(rev.ContentJSON)
	version, err := applyDocContent(l.ctx, l.svcCtx, doc, space, rev.ContentJSON, contentText, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	return &doc_rpc.RollbackRevisionRes{Version: int32(version)}, nil
}
