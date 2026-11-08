package logic

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
)

// randHex 业务 ID 随机段，与 knowledge/platform 同法（crypto/rand）。
func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// genID 生成 prefix_随机段 的业务 ID。
func genID(prefix string) string {
	return prefix + "_" + randHex(16)
}

// trim 去首尾空白。
func trim(s string) string {
	return strings.TrimSpace(s)
}

// sha256hex 正文指纹，用于版本判重。
func sha256hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// rfc3339 把 CustomTime 转 RFC3339 字符串。
func rfc3339(t interface{}) string {
	switch v := t.(type) {
	case time.Time:
		return v.Format(time.RFC3339)
	default:
		return ""
	}
}

// ---- 角色序：用于权限比较与继承。owner>manager>editor>commenter>viewer>无 ----

const (
	roleOwner     = "owner"
	roleManager   = "manager"
	roleEditor    = "editor"
	roleCommenter = "commenter"
	roleViewer    = "viewer"
)

func roleRank(role string) int {
	switch role {
	case roleOwner:
		return 5
	case roleManager:
		return 4
	case roleEditor:
		return 3
	case roleCommenter:
		return 2
	case roleViewer:
		return 1
	default:
		return 0
	}
}

// maxRole 取权限更高者（序大者）。
func maxRole(a, b string) string {
	if roleRank(b) > roleRank(a) {
		return b
	}
	return a
}

// canManage 是否可管理（分享/删改）该资源：owner 或 manager。
func canManage(role string) bool {
	return roleRank(role) >= roleRank(roleManager)
}

// canEdit 是否可编辑正文：owner/manager/editor。
func canEdit(role string) bool {
	return roleRank(role) >= roleRank(roleEditor)
}

// ---- 模型 -> proto item ----

func toSpaceItem(row *doc_models.DocSpace) *doc_rpc.SpaceItem {
	if row == nil {
		return nil
	}
	return &doc_rpc.SpaceItem{
		Id:          uint64(row.Id),
		SpaceId:     row.SpaceID,
		Name:        row.Name,
		Kind:        row.Kind,
		OwnerUserId: row.OwnerUserID,
		ParentId:    row.ParentID,
		Icon:        row.Icon,
		Description: row.Description,
		Status:      int32(row.Status),
		Sort:        int32(row.Sort),
		CreatedAt:   time.Time(row.CreatedAt).Format(time.RFC3339),
	}
}

func toNodeItem(row *doc_models.DocNode) *doc_rpc.NodeItem {
	if row == nil {
		return nil
	}
	return &doc_rpc.NodeItem{
		Id:       uint64(row.Id),
		NodeId:   row.NodeID,
		SpaceId:  row.SpaceID,
		ParentId: row.ParentID,
		Title:    row.Title,
		Icon:     row.Icon,
		Type:     row.Type,
		Sort:     int32(row.Sort),
		DocId:    row.DocID,
	}
}

func toDocItem(row *doc_models.DocDocument) *doc_rpc.DocItem {
	if row == nil {
		return nil
	}
	return &doc_rpc.DocItem{
		Id:           uint64(row.Id),
		DocId:        row.DocID,
		NodeId:       row.NodeID,
		SpaceId:      row.SpaceID,
		Title:        row.Title,
		ContentJson:  row.ContentJSON,
		Version:      int32(row.Version),
		OwnerUserId:  row.OwnerUserID,
		LastEditedBy: row.LastEditedBy,
		Status:       int32(row.Status),
		UpdatedAt:    time.Time(row.UpdatedAt).Format(time.RFC3339),
		CreatedAt:    time.Time(row.CreatedAt).Format(time.RFC3339),
	}
}

func toPermItem(row *doc_models.DocPermission) *doc_rpc.PermItem {
	if row == nil {
		return nil
	}
	return &doc_rpc.PermItem{
		Id:            uint64(row.Id),
		ResourceType:  row.ResourceType,
		ResourceId:    row.ResourceID,
		PrincipalType: row.PrincipalType,
		PrincipalId:   row.PrincipalID,
		Role:          row.Role,
	}
}

func toRevisionItem(row *doc_models.DocRevision) *doc_rpc.RevisionItem {
	if row == nil {
		return nil
	}
	return &doc_rpc.RevisionItem{
		RevisionId: row.RevisionID,
		DocId:      row.DocID,
		Version:    int32(row.Version),
		Title:      row.Title,
		Label:      row.Label,
		CreatedBy:  row.CreatedBy,
		CreatedAt:  time.Time(row.CreatedAt).Format(time.RFC3339),
	}
}
