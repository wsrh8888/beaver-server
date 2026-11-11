package logic

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// nowStr 统一时间字符串，所有对外时间字段都是 RFC3339。
func nowStr() string { return time.Now().Format(time.RFC3339) }

// sentinelDeletedAt 是业务表软删字段的「未删除」值，新建行时显式写入，
// 不依赖数据库 default，避免不同驱动下零值行为不一致。
var sentinelDeletedAt = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

// linkExpired 判断分享链接的 RFC3339 到期时间是否已过；空串表示永不过期。
func linkExpired(expiresAt string) bool {
	if trim(expiresAt) == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return false
	}
	return time.Now().After(t)
}

// audit 在业务库写一条审计。审计与业务数据同库；写失败不回滚业务，但会由上层逻辑记录。
func audit(db *gorm.DB, actor, action, resourceType, resourceID, detail string) {
	_ = db.Create(&doc_models.DocAuditLog{
		LogID:        genID("log"),
		ActorUserID:  actor,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Detail:       detail,
		CreatedAtStr: nowStr(),
	}).Error
}

// auditJSON 用结构化上下文写审计，避免把正文塞进审计表。
func auditJSON(db *gorm.DB, actor, action, resourceType, resourceID string, payload map[string]any) {
	b, _ := json.Marshal(payload)
	audit(db, actor, action, resourceType, resourceID, string(b))
}

// notify 给单个用户写一条文档域通知；自己触发自己的事件不产生通知。
func notify(db *gorm.DB, userID, docID, kind, actor, title, excerpt string) {
	if trim(userID) == "" || userID == actor {
		return
	}
	_ = db.Create(&doc_models.DocNotification{
		NotificationID: genID("ntf"),
		UserID:         userID,
		DocID:          docID,
		Kind:           kind,
		ActorUserID:    actor,
		Title:          title,
		Excerpt:        firstRunes(excerpt, 180),
		CreatedAtStr:   nowStr(),
	}).Error
}

func firstRunes(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}

// splitMentions 解析 "u1,u2" 形式的 @ 列表，去重去空。
func splitMentions(s string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, part := range strings.Split(s, ",") {
		id := trim(part)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}

// enqueueIndex 在正文事务内写一条 pending outbox。平台知识库按 (doc_id, version) 拉
// GetDocContent 建索引，doc 域因此不需要反向依赖 platform_rpc。
func enqueueIndex(tx *gorm.DB, docID string, version int, digest string) {
	_ = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&doc_models.DocIndexOutbox{
		DocID:      docID,
		Version:    version,
		Digest:     digest,
		Status:     "pending",
		CreatedStr: nowStr(),
	}).Error
}

// requireDoc 读取未删文档并校验调用者角色，返回文档行与有效角色。
func requireDoc(db *gorm.DB, docID, userID string, minRole string) (*doc_models.DocDocument, string, error) {
	doc, err := loadDoc(db, docID)
	if err != nil {
		return nil, "", errNotFound
	}
	role, err := resolveDocRole(db, doc.SpaceID, doc.DocID, userID)
	if err != nil {
		return nil, "", err
	}
	if roleRank(role) < roleRank(minRole) {
		return nil, "", errNoPermission
	}
	return doc, role, nil
}

func requireSpace(db *gorm.DB, spaceID, userID string, minRole string) (*doc_models.DocSpace, string, error) {
	space, err := loadSpace(db, spaceID)
	if err != nil {
		return nil, "", errNotFound
	}
	role, err := resolveSpaceRole(db, space, userID)
	if err != nil {
		return nil, "", err
	}
	if roleRank(role) < roleRank(minRole) {
		return nil, "", errNoPermission
	}
	return space, role, nil
}

// pageWindow 统一分页，避免每个端点重复判空。
func pageWindow(page, size int32) (int, int) {
	return normalizePage(page, size)
}

// addTrash 登记一条回收站记录；同一资源重复删除只刷新该行，保证 trash_id 稳定。
func addTrash(tx *gorm.DB, resourceType, resourceID, title, spaceID, parentID, operator, payloadJSON string) {
	var row doc_models.DocTrash
	if err := tx.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).First(&row).Error; err == nil {
		_ = tx.Model(&doc_models.DocTrash{}).Where("trash_id = ?", row.TrashID).Updates(map[string]any{
			"title": title, "space_id": spaceID, "parent_id": parentID,
			"deleted_by": operator, "deleted_at": nowStr(), "payload_json": payloadJSON,
		}).Error
		return
	}
	_ = tx.Create(&doc_models.DocTrash{
		TrashID:      genID("trh"),
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Title:        title,
		SpaceID:      spaceID,
		ParentID:     parentID,
		DeletedBy:    operator,
		DeletedAt:    nowStr(),
		PayloadJSON:  payloadJSON,
	}).Error
}

// trashPayload 把子树节点/文档 ID 序列化为回收站载荷。
func trashPayload(nodeIDs, docIDs []string) string {
	b, _ := json.Marshal(map[string]any{"nodeIds": nodeIDs, "docIds": docIDs})
	return string(b)
}

// zeroDeletedAt 是业务表软删的「未删除」哨兵值，与 models 定义保持一致。
const zeroDeletedAt = "1970-01-01 00:00:00"

// visibleDocQuery 构造「该用户可见的未删文档」查询：
// 自己创建的 + 被显式授权(doc) 的 + 位于自己所属/被授权空间内的。
// 所有检索/列表端点共用，避免各处权限口径不一致造成越权读取。
func visibleDocQuery(db *gorm.DB, userID string) *gorm.DB {
	q := db.Model(&doc_models.DocDocument{}).Where("deleted_at = ?", zeroDeletedAt)
	uid := trim(userID)
	if uid == "" {
		return q.Where("1 = 0")
	}
	var docIDs, spaceIDs []string
	_ = db.Model(&doc_models.DocPermission{}).
		Where("resource_type = ? AND principal_type = ? AND principal_id = ?", "doc", "user", uid).
		Pluck("resource_id", &docIDs).Error
	_ = db.Model(&doc_models.DocPermission{}).
		Where("resource_type = ? AND principal_type = ? AND principal_id = ?", "space", "user", uid).
		Pluck("resource_id", &spaceIDs).Error
	var ownedSpaces []string
	_ = db.Model(&doc_models.DocSpace{}).Where("owner_user_id = ?", uid).
		Pluck("space_id", &ownedSpaces).Error
	spaceIDs = append(spaceIDs, ownedSpaces...)
	return q.Where("owner_user_id = ? OR doc_id IN ? OR space_id IN ?",
		uid, append(docIDs, ""), append(spaceIDs, ""))
}

// snippetAround 取命中词附近的正文片段，用于搜索结果预览。
func snippetAround(text, query string, radius int) string {
	r := []rune(text)
	start := 0
	if q := strings.ToLower(trim(query)); q != "" {
		if pos := strings.Index(strings.ToLower(text), q); pos >= 0 {
			start = len([]rune(text[:pos])) - radius/3
		}
	}
	if start < 0 {
		start = 0
	}
	end := start + radius
	if end > len(r) {
		end = len(r)
	}
	prefix, suffix := "", ""
	if start > 0 {
		prefix = "…"
	}
	if end < len(r) {
		suffix = "…"
	}
	return prefix + string(r[start:end]) + suffix
}

func parseTrashPayload(s string) (nodeIDs, docIDs []string) {
	if trim(s) == "" {
		return nil, nil
	}
	var p struct {
		NodeIDs []string `json:"nodeIds"`
		DocIDs  []string `json:"docIds"`
	}
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return nil, nil
	}
	return p.NodeIDs, p.DocIDs
}

// purgeOneTrash 彻底删除一条回收站记录及其对应的业务数据（含所有子表）。
// 物理删除只在 purge/empty 时发生，普通删除始终是软删 + 登记。
func purgeOneTrash(db *gorm.DB, row *doc_models.DocTrash) error {
	return db.Transaction(func(tx *gorm.DB) error {
		nodeIDs, docIDs := parseTrashPayload(row.PayloadJSON)
		switch row.ResourceType {
		case "doc":
			docIDs = append(docIDs, row.ResourceID)
		case "node":
			if len(nodeIDs) == 0 {
				nodeIDs = append(nodeIDs, row.ResourceID)
			}
		case "space":
			// 空间下的节点/文档以删除时快照为准；再按 space_id 兜底清理残余。
			_ = tx.Where("space_id = ?", row.ResourceID).Delete(&doc_models.DocNode{}).Error
			var extra []string
			_ = tx.Model(&doc_models.DocDocument{}).Where("space_id = ?", row.ResourceID).
				Pluck("doc_id", &extra).Error
			docIDs = append(docIDs, extra...)
			_ = tx.Where("space_id = ?", row.ResourceID).Delete(&doc_models.DocSpace{}).Error
		}
		if len(nodeIDs) > 0 {
			_ = tx.Where("node_id IN ?", nodeIDs).Delete(&doc_models.DocNode{}).Error
		}
		if len(docIDs) > 0 {
			for _, m := range []any{
				&doc_models.DocDocument{}, &doc_models.DocRevision{},
				&doc_models.DocComment{}, &doc_models.DocBlock{}, &doc_models.DocYjsUpdate{},
				&doc_models.DocYjsSnapshot{}, &doc_models.DocRecentVisit{},
				&doc_models.DocFavorite{}, &doc_models.DocIndexOutbox{}, &doc_models.DocNotification{},
			} {
				_ = tx.Where("doc_id IN ?", docIDs).Delete(m).Error
			}
			// 权限与分享链接按 resource_type + resource_id 建模。
			_ = tx.Where("resource_type = ? AND resource_id IN ?", "doc", docIDs).
				Delete(&doc_models.DocPermission{}).Error
			_ = tx.Where("resource_type = ? AND resource_id IN ?", "doc", docIDs).
				Delete(&doc_models.DocShareLink{}).Error
			_ = tx.Where("resource_type = ? AND resource_id IN ?", "doc", docIDs).
				Delete(&doc_models.DocTrash{}).Error
		}
		return tx.Where("trash_id = ?", row.TrashID).Delete(&doc_models.DocTrash{}).Error
	})
}

func toCommentItem(row *doc_models.DocComment, replies []*doc_rpc.CommentItem) *doc_rpc.CommentItem {
	if row == nil {
		return nil
	}
	return &doc_rpc.CommentItem{
		CommentId:  row.CommentID,
		DocId:      row.DocID,
		Epoch:      int32(row.Epoch),
		ParentId:   row.ParentID,
		AnchorId:   row.AnchorID,
		Quote:      row.Quote,
		Content:    row.Content,
		Mentions:   row.Mentions,
		CreatedBy:  row.CreatedBy,
		Resolved:   row.Resolved,
		ResolvedBy: row.ResolvedBy,
		ResolvedAt: row.ResolvedAt,
		CreatedAt:  time.Time(row.CreatedAt).Format(time.RFC3339),
		Replies:    replies,
	}
}

func toShareLinkItem(row *doc_models.DocShareLink) *doc_rpc.ShareLinkItem {
	if row == nil {
		return nil
	}
	return &doc_rpc.ShareLinkItem{
		LinkId:        row.LinkID,
		Token:         row.Token,
		ResourceType:  row.ResourceType,
		ResourceId:    row.ResourceID,
		Scope:         row.Scope,
		Role:          row.Role,
		HasPassword:   row.PasswordHash != "",
		ExpiresAt:     row.ExpiresAt,
		AllowDownload: row.AllowDownload,
		AllowCopy:     row.AllowCopy,
		AllowComment:  row.AllowComment,
		Watermark:     row.Watermark,
		ViewCount:     row.ViewCount,
		CreatedBy:     row.CreatedBy,
		Status:        int32(row.Status),
		Url:           "/doc/s/" + row.Token,
	}
}

func toTrashItem(row *doc_models.DocTrash) *doc_rpc.TrashItem {
	if row == nil {
		return nil
	}
	return &doc_rpc.TrashItem{
		TrashId:      row.TrashID,
		ResourceType: row.ResourceType,
		ResourceId:   row.ResourceID,
		Title:        row.Title,
		SpaceId:      row.SpaceID,
		ParentId:     row.ParentID,
		DeletedBy:    row.DeletedBy,
		DeletedAt:    row.DeletedAt,
	}
}

func toTemplateItem(row *doc_models.DocTemplate) *doc_rpc.TemplateItem {
	if row == nil {
		return nil
	}
	return &doc_rpc.TemplateItem{
		TemplateId:  row.TemplateID,
		Name:        row.Name,
		Category:    row.Category,
		Description: row.Description,
		Cover:       row.Cover,
		ContentJson: row.ContentJSON,
		DocType:     row.DocType,
		Scope:       row.Scope,
		OwnerUserId: row.OwnerUserID,
		UseCount:    int32(row.UseCount),
		Status:      int32(row.Status),
		CreatedAt:   time.Time(row.CreatedAt).Format(time.RFC3339),
	}
}

func toDocBlockItem(row *doc_models.DocBlock) *doc_rpc.DocBlockItem {
	if row == nil {
		return nil
	}
	return &doc_rpc.DocBlockItem{
		BlockId:     row.BlockID,
		DocId:       row.DocID,
		Kind:        row.Kind,
		Title:       row.Title,
		PayloadJson: row.PayloadJSON,
		Version:     int32(row.Version),
		UpdatedBy:   row.UpdatedBy,
		UpdatedAt:   time.Time(row.UpdatedAt).Format(time.RFC3339),
	}
}

func toNotificationItem(row *doc_models.DocNotification) *doc_rpc.NotificationItem {
	if row == nil {
		return nil
	}
	return &doc_rpc.NotificationItem{
		NotificationId: row.NotificationID,
		UserId:         row.UserID,
		DocId:          row.DocID,
		Kind:           row.Kind,
		ActorUserId:    row.ActorUserID,
		Title:          row.Title,
		Excerpt:        row.Excerpt,
		IsRead:         row.IsRead,
		CreatedAt:      row.CreatedAtStr,
	}
}

func toAuditLogItem(row *doc_models.DocAuditLog) *doc_rpc.AuditLogItem {
	if row == nil {
		return nil
	}
	return &doc_rpc.AuditLogItem{
		LogId:        row.LogID,
		ActorUserId:  row.ActorUserID,
		Action:       row.Action,
		ResourceType: row.ResourceType,
		ResourceId:   row.ResourceID,
		Detail:       row.Detail,
		CreatedAt:    row.CreatedAtStr,
	}
}

// ctxOf 保留 logic 的 context，供领域函数统一使用。
func ctxOf(ctx context.Context, s *svc.ServiceContext) (*gorm.DB, context.Context) {
	_ = ctx
	return s.DB, ctx
}
