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

type DeleteDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDocLogic {
	return &DeleteDocLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DeleteDocLogic", ctx),
	}
}

func (l *DeleteDocLogic) DeleteDoc(in *doc_rpc.DeleteDocReq) (*doc_rpc.DeleteDocRes, error) {
	doc, err := loadDoc(l.svcCtx.DB, in.DocId)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveDocRole(l.svcCtx.DB, doc.SpaceID, doc.DocID, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canManage(role) {
		return nil, errNoPermission
	}

	now := time.Now()
	if err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&doc_models.DocDocument{}).Where("doc_id = ?", doc.DocID).
			Update("deleted_at", now).Error; err != nil {
			return err
		}
		// 关联节点一并软删。
		if doc.NodeID != "" {
			if err := tx.Model(&doc_models.DocNode{}).Where("node_id = ?", doc.NodeID).
				Update("deleted_at", now).Error; err != nil {
				return err
			}
		}
		addTrash(tx, "doc", doc.DocID, doc.Title, doc.SpaceID, doc.NodeID, in.OperatorUserId, "")
		auditJSON(tx, in.OperatorUserId, "doc.delete", "doc", doc.DocID,
			map[string]any{"title": doc.Title, "spaceId": doc.SpaceID})
		return nil
	}); err != nil {
		l.logger.Error(model.LogMsg{Text: "删除文档失败", Data: err})
		return nil, err
	}

	return &doc_rpc.DeleteDocRes{}, nil
}
