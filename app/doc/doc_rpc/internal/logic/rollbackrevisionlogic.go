package logic

import (
	"context"
	"strconv"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"

	"gorm.io/gorm"
)

type RollbackRevisionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRollbackRevisionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RollbackRevisionLogic {
	return &RollbackRevisionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("RollbackRevisionLogic", ctx),
	}
}

func (l *RollbackRevisionLogic) RollbackRevision(in *doc_rpc.RollbackRevisionReq) (*doc_rpc.RollbackRevisionRes, error) {
	doc, err := loadDoc(l.svcCtx.DB, in.DocId)
	if err != nil {
		return nil, errNotFound
	}
	space, err := loadSpace(l.svcCtx.DB, doc.SpaceID)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveDocRole(l.svcCtx.DB, doc.SpaceID, doc.DocID, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canEdit(role) {
		return nil, errNoPermission
	}

	var rev doc_models.DocRevision
	if err := l.svcCtx.DB.Where("revision_id = ? AND doc_id = ?", in.RevisionId, doc.DocID).First(&rev).Error; err != nil {
		return nil, errNotFound
	}

	// 回滚 = 用该版本正文生成新版本（保留历史线性，不删中间版本）。
	contentText := extractTextFromJSON(rev.ContentJSON)
	version, err := applyDocContent(l.ctx, l.svcCtx, doc, space, rev.ContentJSON, contentText, in.OperatorUserId, "revision.rollback")
	if err != nil {
		return nil, err
	}

	// CRDT 状态无法「倒退」：递增 epoch，房间名由 docId 变为 docId#epoch，
	// 旧纪元的 update 日志与快照整体作废，所有客户端以新纪元重新同步到回滚后的正文。
	newEpoch := doc.Epoch + 1
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&doc_models.DocDocument{}).Where("doc_id = ?", doc.DocID).
			Update("epoch", newEpoch).Error; err != nil {
			return err
		}
		if err := tx.Where("doc_id = ?", doc.DocID).Delete(&doc_models.DocYjsUpdate{}).Error; err != nil {
			return err
		}
		return tx.Where("doc_id = ?", doc.DocID).Delete(&doc_models.DocYjsSnapshot{}).Error
	})
	if err != nil {
		return nil, err
	}
	doc.Epoch = newEpoch

	auditJSON(l.svcCtx.DB, in.OperatorUserId, "room.epoch_bump", "doc", doc.DocID, map[string]any{
		"epoch": newEpoch, "version": version, "revisionId": in.RevisionId,
	})
	notify(l.svcCtx.DB, doc.OwnerUserID, doc.DocID, "rollback", in.OperatorUserId,
		"《"+doc.Title+"》已回滚到版本 "+strconv.Itoa(rev.Version), contentText)

	return &doc_rpc.RollbackRevisionRes{Version: int32(version), Epoch: int32(newEpoch)}, nil
}
