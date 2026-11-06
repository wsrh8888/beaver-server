package logic

import (
	"context"

	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	knowledgecli "beaver/app/knowledge/knowledge_rpc/knowledge"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ReindexDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewReindexDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReindexDocLogic {
	return &ReindexDocLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ReindexDocLogic", ctx),
	}
}

// ReindexDoc 触发文档 RAG 重新索引：投影缺失时先补投影（projectDoc 会创建并回填 kid），
// 再调用 knowledge_rpc.ReindexDocument 让文档重新入队分段->embedding->OpenSearch。
func (l *ReindexDocLogic) ReindexDoc(in *doc_rpc.ReindexDocReq) (*doc_rpc.ReindexDocRes, error) {
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

	if !l.svcCtx.Config.Rag.Enabled {
		return &doc_rpc.ReindexDocRes{Ok: false}, nil
	}

	// 补投影：无 knowledge_doc_id 时 projectDoc 会创建对应 knowledge 文档并回填。
	projectDoc(l.ctx, l.svcCtx, space, doc)
	kid := doc.KnowledgeDocID
	if kid == "" {
		return &doc_rpc.ReindexDocRes{Ok: false}, nil
	}

	if _, err := l.svcCtx.KnowledgeRpc.ReindexDocument(l.ctx, &knowledgecli.ReindexDocumentReq{DocumentId: kid}); err != nil {
		l.logger.Error(model.LogMsg{Text: "触发文档重索引失败", Data: map[string]any{"docId": doc.DocID, "kid": kid, "err": err.Error()}})
		return &doc_rpc.ReindexDocRes{Ok: false, KnowledgeDocId: kid}, nil
	}
	return &doc_rpc.ReindexDocRes{Ok: true, KnowledgeDocId: kid}, nil
}
