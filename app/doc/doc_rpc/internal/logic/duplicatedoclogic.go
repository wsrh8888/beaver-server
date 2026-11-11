package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"

	"gorm.io/gorm"
)

type DuplicateDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDuplicateDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DuplicateDocLogic {
	return &DuplicateDocLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DuplicateDocLogic", ctx),
	}
}

func (l *DuplicateDocLogic) DuplicateDoc(in *doc_rpc.DuplicateDocReq) (*doc_rpc.DuplicateDocRes, error) {
	src, _, err := requireDoc(l.svcCtx.DB, in.DocId, in.OperatorUserId, roleViewer)
	if err != nil {
		return nil, err
	}
	space, _, err := requireSpace(l.svcCtx.DB, in.SpaceId, in.OperatorUserId, roleEditor)
	if err != nil {
		return nil, err
	}
	actor := trim(in.OperatorUserId)
	title := trim(in.Title)
	if title == "" {
		title = src.Title + " 副本"
	}
	nodeID, docID := genID("node"), genID("doc")
	row := doc_models.DocDocument{
		DocID:         docID,
		NodeID:        nodeID,
		SpaceID:       space.SpaceID,
		Title:         title,
		ContentJSON:   src.ContentJSON,
		ContentText:   src.ContentText,
		ContentDigest: src.ContentDigest,
		Version:       1,
		Epoch:         1,
		DocType:       src.DocType,
		Icon:          src.Icon,
		Cover:         src.Cover,
		RagEnabled:    src.RagEnabled,
		OwnerUserID:   actor,
		Status:        1,
		DeletedAt:     sentinelDeletedAt,
	}
	node := doc_models.DocNode{
		NodeID: nodeID, SpaceID: space.SpaceID, ParentID: trim(in.ParentNodeId),
		Title: title, Type: "page", Icon: src.Icon, DocID: docID,
	}
	var blocks []doc_models.DocBlock
	_ = l.svcCtx.DB.Where("doc_id = ?", src.DocID).Find(&blocks).Error

	if err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&node).Error; err != nil {
			return err
		}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
		for i := range blocks {
			cp := blocks[i]
			cp.Id = 0
			cp.BlockID = genID("blk")
			cp.DocID = docID
			cp.Version = 1
			cp.UpdatedBy = actor
			if err := tx.Create(&cp).Error; err != nil {
				return err
			}
		}
		if row.RagEnabled {
			enqueueIndex(tx, docID, row.Version, row.ContentDigest)
		}
		auditJSON(tx, actor, "doc.duplicate", "doc", docID,
			map[string]any{"sourceDocId": src.DocID, "spaceId": space.SpaceID})
		return nil
	}); err != nil {
		return nil, err
	}
	return &doc_rpc.DuplicateDocRes{DocId: docID, NodeId: nodeID}, nil
}
