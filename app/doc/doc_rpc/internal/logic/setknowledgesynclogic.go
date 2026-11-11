package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type SetKnowledgeSyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSetKnowledgeSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetKnowledgeSyncLogic {
	return &SetKnowledgeSyncLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("SetKnowledgeSyncLogic", ctx),
	}
}

func (l *SetKnowledgeSyncLogic) SetKnowledgeSync(in *doc_rpc.SetKnowledgeSyncReq) (*doc_rpc.SetKnowledgeSyncRes, error) {
	doc, _, err := requireDoc(l.svcCtx.DB, in.DocId, in.OperatorUserId, roleManager)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.DB.Model(&doc_models.DocDocument{}).
		Where("doc_id = ?", doc.DocID).Update("rag_enabled", in.RagEnabled).Error; err != nil {
		return nil, err
	}
	// 首次开启时把当前正文补入索引队列。
	if in.RagEnabled && !doc.RagEnabled {
		enqueueIndex(l.svcCtx.DB, doc.DocID, doc.Version, doc.ContentDigest)
	}
	auditJSON(l.svcCtx.DB, in.OperatorUserId, "doc.knowledge_sync", "doc", doc.DocID,
		map[string]any{"ragEnabled": in.RagEnabled})
	return &doc_rpc.SetKnowledgeSyncRes{RagEnabled: in.RagEnabled}, nil
}
