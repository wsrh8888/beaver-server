package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"

	"gorm.io/gorm/clause"
)

type ReindexDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewReindexDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReindexDocLogic {
	return &ReindexDocLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ReindexDocLogic", ctx),
	}
}

func (l *ReindexDocLogic) ReindexDoc(in *doc_rpc.ReindexDocReq) (*doc_rpc.ReindexDocRes, error) {
	doc, _, err := requireDoc(l.svcCtx.DB, in.DocId, in.OperatorUserId, roleEditor)
	if err != nil {
		return nil, err
	}
	if !doc.RagEnabled {
		return &doc_rpc.ReindexDocRes{
			Version: int32(doc.Version), Status: "disabled", Digest: doc.ContentDigest,
		}, nil
	}
	// 强制重建：同一 (doc_id, version) 若已是 done，也要重置为 pending 并清空重试次数。
	row := doc_models.DocIndexOutbox{
		DocID:      doc.DocID,
		Version:    doc.Version,
		Digest:     doc.ContentDigest,
		Status:     "pending",
		Attempts:   0,
		CreatedStr: nowStr(),
	}
	if err := l.svcCtx.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "doc_id"}, {Name: "version"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"digest", "status", "attempts", "created_str",
		}),
	}).Create(&row).Error; err != nil {
		return nil, err
	}
	auditJSON(l.svcCtx.DB, in.OperatorUserId, "doc.reindex", "doc", doc.DocID,
		map[string]any{"version": doc.Version})
	return &doc_rpc.ReindexDocRes{
		Version: int32(doc.Version), Status: "queued", Digest: doc.ContentDigest,
	}, nil
}
