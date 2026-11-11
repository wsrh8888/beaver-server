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

type DeleteSpaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSpaceLogic {
	return &DeleteSpaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DeleteSpaceLogic", ctx),
	}
}

func (l *DeleteSpaceLogic) DeleteSpace(in *doc_rpc.DeleteSpaceReq) (*doc_rpc.DeleteSpaceRes, error) {
	space, err := loadSpace(l.svcCtx.DB, in.SpaceId)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveSpaceRole(l.svcCtx.DB, space, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canManage(role) {
		return nil, errNoPermission
	}

	now := time.Now()
	var nodeIDs, docIDs []string
	if err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&doc_models.DocSpace{}).Where("space_id = ?", space.SpaceID).
			Update("deleted_at", now).Error; err != nil {
			return err
		}
		var nodes []doc_models.DocNode
		if err := tx.Where("space_id = ? AND deleted_at = ?", space.SpaceID, zeroDeletedAt).
			Find(&nodes).Error; err != nil {
			return err
		}
		for i := range nodes {
			nodeIDs = append(nodeIDs, nodes[i].NodeID)
		}
		if len(nodeIDs) > 0 {
			if err := tx.Model(&doc_models.DocNode{}).Where("node_id IN ?", nodeIDs).
				Update("deleted_at", now).Error; err != nil {
				return err
			}
		}
		var docs []doc_models.DocDocument
		if err := tx.Where("space_id = ? AND deleted_at = ?", space.SpaceID, zeroDeletedAt).
			Find(&docs).Error; err != nil {
			return err
		}
		for i := range docs {
			docIDs = append(docIDs, docs[i].DocID)
		}
		if len(docIDs) > 0 {
			if err := tx.Model(&doc_models.DocDocument{}).Where("doc_id IN ?", docIDs).
				Update("deleted_at", now).Error; err != nil {
				return err
			}
		}
		addTrash(tx, "space", space.SpaceID, space.Name, space.SpaceID, space.ParentID,
			in.OperatorUserId, trashPayload(nodeIDs, docIDs))
		auditJSON(tx, in.OperatorUserId, "space.delete", "space", space.SpaceID,
			map[string]any{"name": space.Name, "nodes": len(nodeIDs), "docs": len(docIDs)})
		return nil
	}); err != nil {
		l.logger.Error(model.LogMsg{Text: "删除空间失败", Data: err})
		return nil, err
	}
	return &doc_rpc.DeleteSpaceRes{}, nil
}
