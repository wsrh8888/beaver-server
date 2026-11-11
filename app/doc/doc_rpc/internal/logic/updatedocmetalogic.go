package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type UpdateDocMetaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateDocMetaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDocMetaLogic {
	return &UpdateDocMetaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("UpdateDocMetaLogic", ctx),
	}
}

func (l *UpdateDocMetaLogic) UpdateDocMeta(in *doc_rpc.UpdateDocMetaReq) (*doc_rpc.UpdateDocMetaRes, error) {
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

	upd := map[string]any{}
	if title := trim(in.Title); title != "" {
		upd["title"] = title
	}
	if in.Status > 0 {
		upd["status"] = int(in.Status)
	}
	if icon := trim(in.Icon); icon != "" {
		upd["icon"] = icon
	}
	if cover := trim(in.Cover); cover != "" {
		upd["cover"] = cover
	}
	if dt := trim(in.DocType); dt != "" {
		upd["doc_type"] = dt
	}
	// bool 的 false 与「未设置」无法区分，用 rag_enabled_set 显式表达意图。
	if in.RagEnabledSet {
		upd["rag_enabled"] = in.RagEnabled
	}
	if len(upd) == 0 {
		return &doc_rpc.UpdateDocMetaRes{}, nil
	}

	if err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&doc_models.DocDocument{}).Where("doc_id = ?", doc.DocID).Updates(upd).Error; err != nil {
			return err
		}
		// 首次打开知识库同步时补一条索引任务，让已有正文立刻可被检索。
		if in.RagEnabledSet && in.RagEnabled && !doc.RagEnabled {
			enqueueIndex(tx, doc.DocID, doc.Version, doc.ContentDigest)
		}
		auditJSON(tx, in.OperatorUserId, "doc.update_meta", "doc", doc.DocID, upd)
		return nil
	}); err != nil {
		l.logger.Error(model.LogMsg{Text: "更新文档元信息失败", Data: err})
		return nil, err
	}

	if title := trim(in.Title); title != "" {
		doc.Title = title
	}
	if in.Status > 0 {
		doc.Status = int(in.Status)
	}
	if in.RagEnabledSet {
		doc.RagEnabled = in.RagEnabled
	}

	return &doc_rpc.UpdateDocMetaRes{}, nil
}
