package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"

	"gorm.io/gorm"
)

type RestoreTrashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRestoreTrashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreTrashLogic {
	return &RestoreTrashLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("RestoreTrashLogic", ctx),
	}
}

func (l *RestoreTrashLogic) RestoreTrash(in *doc_rpc.RestoreTrashReq) (*doc_rpc.RestoreTrashRes, error) {
	actor := trim(in.OperatorUserId)
	var row doc_models.DocTrash
	if err := l.svcCtx.DB.Where("trash_id = ?", in.TrashId).First(&row).Error; err != nil {
		return nil, errNotFound
	}
	if row.DeletedBy != "" && row.DeletedBy != actor {
		return nil, errNoPermission
	}

	nodeIDs, docIDs := parseTrashPayload(row.PayloadJSON)
	switch row.ResourceType {
	case "doc":
		docIDs = append(docIDs, row.ResourceID)
	case "node":
		if len(nodeIDs) == 0 {
			nodeIDs = append(nodeIDs, row.ResourceID)
		}
	}

	if err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if row.ResourceType == "space" {
			if err := tx.Model(&doc_models.DocSpace{}).Where("space_id = ?", row.ResourceID).
				Update("deleted_at", zeroDeletedAt).Error; err != nil {
				return err
			}
		}
		if len(nodeIDs) > 0 {
			if err := tx.Model(&doc_models.DocNode{}).Where("node_id IN ?", nodeIDs).
				Update("deleted_at", zeroDeletedAt).Error; err != nil {
				return err
			}
		}
		if len(docIDs) > 0 {
			if err := tx.Model(&doc_models.DocDocument{}).Where("doc_id IN ?", docIDs).
				Update("deleted_at", zeroDeletedAt).Error; err != nil {
				return err
			}
		}
		return tx.Where("trash_id = ?", row.TrashID).Delete(&doc_models.DocTrash{}).Error
	}); err != nil {
		return nil, err
	}
	auditJSON(l.svcCtx.DB, actor, "trash.restore", row.ResourceType, row.ResourceID,
		map[string]any{"nodes": len(nodeIDs), "docs": len(docIDs)})
	return &doc_rpc.RestoreTrashRes{ResourceType: row.ResourceType, ResourceId: row.ResourceID}, nil
}
