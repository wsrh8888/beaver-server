package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	knowledgecli "beaver/app/knowledge/knowledge_rpc/knowledge"

	"gorm.io/gorm"
)

// principalRolesFor 查询某用户在给定资源集上的显式授权角色（user 主体 + 其 group/dept 由上层扩展，本期仅 user）。
// 返回 resource_id -> role。
func explicitRole(db *gorm.DB, resourceType, resourceID, userID string) (string, error) {
	if trim(userID) == "" {
		return "", nil
	}
	var row doc_models.DocPermission
	err := db.Where("resource_type = ? AND resource_id = ? AND principal_type = ? AND principal_id = ?",
		resourceType, resourceID, "user", userID).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return row.Role, nil
}

// resolveSpaceRole 计算用户对空间的有效角色：显式授权 与 归属人（owner_user_id）取高者。
func resolveSpaceRole(db *gorm.DB, space *doc_models.DocSpace, userID string) (string, error) {
	if space == nil {
		return "", nil
	}
	role := ""
	if trim(userID) != "" && trim(space.OwnerUserID) == trim(userID) {
		role = roleOwner
	}
	exp, err := explicitRole(db, "space", space.SpaceID, userID)
	if err != nil {
		return "", err
	}
	return maxRole(role, exp), nil
}

// resolveDocRole 计算用户对文档的有效角色：文档显式授权 与 所属空间授权 取高者。
func resolveDocRole(db *gorm.DB, spaceID, docID, userID string) (string, error) {
	docRole, err := explicitRole(db, "doc", docID, userID)
	if err != nil {
		return "", err
	}
	spaceRole, err := explicitRole(db, "space", spaceID, userID)
	if err != nil {
		return "", err
	}
	return maxRole(docRole, spaceRole), nil
}

// loadDoc 取未删文档行。
func loadDoc(db *gorm.DB, docID string) (*doc_models.DocDocument, error) {
	var row doc_models.DocDocument
	if err := db.Where("doc_id = ? AND deleted_at = ?", docID, "1970-01-01 00:00:00").First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// loadSpace 取未删空间行。
func loadSpace(db *gorm.DB, spaceID string) (*doc_models.DocSpace, error) {
	var row doc_models.DocSpace
	if err := db.Where("space_id = ? AND deleted_at = ?", spaceID, "1970-01-01 00:00:00").First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// loadNode 取未删节点行。
func loadNode(db *gorm.DB, nodeID string) (*doc_models.DocNode, error) {
	var row doc_models.DocNode
	if err := db.Where("node_id = ? AND deleted_at = ?", nodeID, "1970-01-01 00:00:00").First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ---- RAG 投影：把文档正文单向同步进 knowledge 域（collection=空间，document=文档）----
// 全部尽力而为：投影失败只记录，不阻断文档主流程（知识库是派生层）。

// ensureSpaceCollection 保证空间有对应的 knowledge collection，回填 space.knowledge_collection_id。
func ensureSpaceCollection(ctx context.Context, svcCtx *svc.ServiceContext, space *doc_models.DocSpace) string {
	if trim(space.KnowledgeCollectionID) != "" {
		return space.KnowledgeCollectionID
	}
	if !svcCtx.Config.Rag.Enabled {
		return ""
	}
	aclMode := "internal"
	if space.Kind == "owner" {
		aclMode = "user_scoped"
	}
	kind := svcCtx.Config.Rag.CollectionKind
	if trim(kind) == "" {
		kind = "docs"
	}
	res, err := svcCtx.KnowledgeRpc.CreateCollection(ctx, &knowledgecli.CreateCollectionReq{
		Name:        space.Name,
		Kind:        kind,
		AclMode:     aclMode,
		Description: "doc space projection",
	})
	if err != nil || res == nil {
		return ""
	}
	cid := res.GetCollectionId()
	if cid == "" {
		return ""
	}
	_ = svcCtx.DB.Model(&doc_models.DocSpace{}).Where("space_id = ?", space.SpaceID).
		Update("knowledge_collection_id", cid).Error
	space.KnowledgeCollectionID = cid
	return cid
}

// projectDoc 保存文档正文后投影进 RAG：正文判重（digest）变化才重新入队。
func projectDoc(ctx context.Context, svcCtx *svc.ServiceContext, space *doc_models.DocSpace, doc *doc_models.DocDocument) {
	if !svcCtx.Config.Rag.Enabled || !doc.RagEnabled || space == nil {
		return
	}
	cid := ensureSpaceCollection(ctx, svcCtx, space)
	if cid == "" {
		return
	}
	if trim(doc.KnowledgeDocID) == "" {
		res, err := svcCtx.KnowledgeRpc.CreateDocument(ctx, &knowledgecli.CreateDocumentReq{
			CollectionId: cid,
			Title:        doc.Title,
			DocType:      "doc",
			SourceRef:    "doc:" + doc.DocID,
			ContentText:  doc.ContentText,
			OwnerUserId:  doc.OwnerUserID,
			RagEnabled:   true,
		})
		if err != nil || res == nil {
			return
		}
		kid := res.GetDocumentId()
		if kid != "" {
			_ = svcCtx.DB.Model(&doc_models.DocDocument{}).Where("doc_id = ?", doc.DocID).
				Update("knowledge_doc_id", kid).Error
			doc.KnowledgeDocID = kid
		}
		return
	}
	_, _ = svcCtx.KnowledgeRpc.UpdateDocument(ctx, &knowledgecli.UpdateDocumentReq{
		DocumentId:   doc.KnowledgeDocID,
		Title:        doc.Title,
		ContentText:  doc.ContentText,
		UpdatedBy:    doc.LastEditedBy,
		RagEnabled:   doc.RagEnabled,
		Status:       int32(doc.Status),
	})
}

// deleteDocProjection 文档软删/停用/关闭检索时清掉 RAG 投影行。
func deleteDocProjection(ctx context.Context, svcCtx *svc.ServiceContext, doc *doc_models.DocDocument) {
	if !svcCtx.Config.Rag.Enabled || trim(doc.KnowledgeDocID) == "" {
		return
	}
	_, _ = svcCtx.KnowledgeRpc.DeleteDocument(ctx, &knowledgecli.DeleteDocumentReq{DocumentId: doc.KnowledgeDocID})
}

// deleteSpaceProjection 空间删除时整库清投影（DeleteCollection 级联软删其下文档并清索引）。
func deleteSpaceProjection(ctx context.Context, svcCtx *svc.ServiceContext, space *doc_models.DocSpace) {
	if !svcCtx.Config.Rag.Enabled || trim(space.KnowledgeCollectionID) == "" {
		return
	}
	_, _ = svcCtx.KnowledgeRpc.DeleteCollection(ctx, &knowledgecli.DeleteCollectionReq{CollectionId: space.KnowledgeCollectionID})
}
