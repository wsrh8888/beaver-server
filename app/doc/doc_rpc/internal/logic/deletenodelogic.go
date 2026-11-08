package logic

import (
	"context"
	"time"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type DeleteNodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNodeLogic {
	return &DeleteNodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DeleteNodeLogic", ctx),
	}
}

func (l *DeleteNodeLogic) DeleteNode(in *doc_rpc.DeleteNodeReq) (*doc_rpc.DeleteNodeRes, error) {
	node, err := loadNode(l.svcCtx.DB, in.NodeId)
	if err != nil {
		return nil, errNotFound
	}
	space, err := loadSpace(l.svcCtx.DB, node.SpaceID)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveSpaceRole(l.svcCtx.DB, space, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canEdit(role) {
		return nil, errNoPermission
	}

	// 递归收集子树节点（含自身）。
	nodeIDs, docIDs := collectSubtree(l.svcCtx.DB, node.NodeID)
	now := time.Now()

	if err := l.svcCtx.DB.Model(&doc_models.DocNode{}).Where("node_id IN ?", nodeIDs).Update("deleted_at", now).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "删除节点失败", Data: err})
		return nil, err
	}
	if len(docIDs) > 0 {
		if err := l.svcCtx.DB.Model(&doc_models.DocDocument{}).Where("doc_id IN ?", docIDs).Update("deleted_at", now).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "删除节点文档失败", Data: err})
		}
	}

	return &doc_rpc.DeleteNodeRes{}, nil
}

// collectSubtree 广度优先收集 node 及其所有后代的 node_id 与关联 doc_id。
func collectSubtree(db *gorm.DB, rootID string) (nodeIDs, docIDs []string) {
	visited := map[string]bool{}
	queue := []string{rootID}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if visited[cur] {
			continue
		}
		visited[cur] = true
		nodeIDs = append(nodeIDs, cur)

		var node doc_models.DocNode
		if err := db.Where("node_id = ? AND deleted_at = ?", cur, "1970-01-01 00:00:00").First(&node).Error; err != nil {
			continue
		}
		if node.DocID != "" {
			docIDs = append(docIDs, node.DocID)
		}

		var children []doc_models.DocNode
		db.Where("parent_id = ? AND deleted_at = ?", cur, "1970-01-01 00:00:00").Find(&children)
		for _, c := range children {
			queue = append(queue, c.NodeID)
		}
	}
	return nodeIDs, docIDs
}
