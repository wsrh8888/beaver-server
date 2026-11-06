package logic

import (
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
	"beaver/app/doc/doc_api/internal/types"
)

// convSpaceItem 把 doc_rpc 的 item 映射为 api 出参（去掉内部自增 id，字段对齐 proto）。
func convSpaceItem(it *doc_rpc.SpaceItem) types.SpaceItem {
	if it == nil {
		return types.SpaceItem{}
	}
	return types.SpaceItem{
		SpaceId:     it.SpaceId,
		Name:        it.Name,
		Kind:        it.Kind,
		OwnerUserId: it.OwnerUserId,
		ParentId:    it.ParentId,
		Icon:        it.Icon,
		Description: it.Description,
		Status:      int(it.Status),
		Sort:        int(it.Sort),
		CreatedAt:   it.CreatedAt,
	}
}

func convSpaceItems(list []*doc_rpc.SpaceItem) []types.SpaceItem {
	out := make([]types.SpaceItem, 0, len(list))
	for _, it := range list {
		out = append(out, convSpaceItem(it))
	}
	return out
}

func convNodeItem(it *doc_rpc.NodeItem) types.NodeItem {
	if it == nil {
		return types.NodeItem{}
	}
	return types.NodeItem{
		NodeId:   it.NodeId,
		SpaceId:  it.SpaceId,
		ParentId: it.ParentId,
		Title:    it.Title,
		Icon:     it.Icon,
		Type:     it.Type,
		Sort:     int(it.Sort),
		DocId:    it.DocId,
	}
}

func convNodeItems(list []*doc_rpc.NodeItem) []types.NodeItem {
	out := make([]types.NodeItem, 0, len(list))
	for _, it := range list {
		out = append(out, convNodeItem(it))
	}
	return out
}

func convDocItem(it *doc_rpc.DocItem) *types.DocItem {
	if it == nil {
		return nil
	}
	return &types.DocItem{
		DocId:        it.DocId,
		NodeId:       it.NodeId,
		SpaceId:      it.SpaceId,
		Title:        it.Title,
		ContentJson:  it.ContentJson,
		Version:      int(it.Version),
		OwnerUserId:  it.OwnerUserId,
		LastEditedBy: it.LastEditedBy,
		RagEnabled:   it.RagEnabled,
		Status:       int(it.Status),
		UpdatedAt:    it.UpdatedAt,
		CreatedAt:    it.CreatedAt,
	}
}

func convPermItem(it *doc_rpc.PermItem) types.PermItem {
	if it == nil {
		return types.PermItem{}
	}
	return types.PermItem{
		ResourceType:  it.ResourceType,
		ResourceId:    it.ResourceId,
		PrincipalType: it.PrincipalType,
		PrincipalId:   it.PrincipalId,
		Role:          it.Role,
	}
}

func convPermItems(list []*doc_rpc.PermItem) []types.PermItem {
	out := make([]types.PermItem, 0, len(list))
	for _, it := range list {
		out = append(out, convPermItem(it))
	}
	return out
}

func convRevisionItem(it *doc_rpc.RevisionItem) types.RevisionItem {
	if it == nil {
		return types.RevisionItem{}
	}
	return types.RevisionItem{
		RevisionId: it.RevisionId,
		DocId:      it.DocId,
		Version:    int(it.Version),
		Title:      it.Title,
		Label:      it.Label,
		CreatedBy:  it.CreatedBy,
		CreatedAt:  it.CreatedAt,
	}
}

func convRevisionItems(list []*doc_rpc.RevisionItem) []types.RevisionItem {
	out := make([]types.RevisionItem, 0, len(list))
	for _, it := range list {
		out = append(out, convRevisionItem(it))
	}
	return out
}
