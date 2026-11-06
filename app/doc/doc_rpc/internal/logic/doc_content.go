package logic

import (
	"context"
	"errors"
	"encoding/json"
	"strings"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"

	"gorm.io/gorm"
)

var (
	errNotFound    = errors.New("doc: 资源不存在")
	errNoPermission = errors.New("doc: 无权限")
	errInvalidArg  = errors.New("doc: 参数非法")
)

// applyDocContent 写入正文：更新 content/digest/version/last_edited_by，生成一条 revision，投影 RAG。
// 由 SaveDocContent / RollbackRevision / CommitDocSnapshot 复用，保证版本语义单一来源。
func applyDocContent(ctx context.Context, svcCtx *svc.ServiceContext, doc *doc_models.DocDocument, space *doc_models.DocSpace, contentJSON, contentText, operator string) (int, error) {
	newVersion := doc.Version + 1
	digest := sha256hex(contentText)

	err := svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		upd := map[string]any{
			"content_json":   contentJSON,
			"content_text":   contentText,
			"content_digest": digest,
			"version":        newVersion,
			"last_edited_by": operator,
		}
		if err := tx.Model(&doc_models.DocDocument{}).Where("doc_id = ?", doc.DocID).Updates(upd).Error; err != nil {
			return err
		}
		rev := doc_models.DocRevision{
			RevisionID:  genID("rev"),
			DocID:       doc.DocID,
			Version:     newVersion,
			Title:       doc.Title,
			ContentJSON: contentJSON,
			CreatedBy:   operator,
		}
		return tx.Create(&rev).Error
	})
	if err != nil {
		return 0, err
	}

	// 回填内存态并投影 RAG。
	doc.ContentJSON = contentJSON
	doc.ContentText = contentText
	doc.ContentDigest = digest
	doc.Version = newVersion
	doc.LastEditedBy = operator
	projectDoc(ctx, svcCtx, space, doc)
	return newVersion, nil
}

// skeletonContentJSON 空文档初始 block（ProseMirror/Tiptap 兼容的单空段落）。
const skeletonContentJSON = `{"type":"doc","content":[{"type":"paragraph"}]}`

// extractTextFromJSON 从 ProseMirror block JSON 递归抽取纯文本（供回滚后重建 RAG 正文）。
// 文本节点取 text 字段；块级节点（段落/标题等）之间以换行分隔。解析失败返回空串。
func extractTextFromJSON(contentJSON string) string {
	if trim(contentJSON) == "" {
		return ""
	}
	var root any
	if err := json.Unmarshal([]byte(contentJSON), &root); err != nil {
		return ""
	}
	var sb strings.Builder
	walkExtractText(root, &sb)
	return strings.TrimSpace(sb.String())
}

func walkExtractText(node any, sb *strings.Builder) {
	switch v := node.(type) {
	case map[string]any:
		if t, ok := v["text"].(string); ok {
			sb.WriteString(t)
		}
		if children, ok := v["content"].([]any); ok {
			for _, c := range children {
				walkExtractText(c, sb)
			}
			// 块级节点结束后补换行，避免段落粘连。
			if typ, _ := v["type"].(string); typ != "text" && typ != "" {
				sb.WriteString("\n")
			}
		}
	case []any:
		for _, c := range v {
			walkExtractText(c, sb)
		}
	}
}
