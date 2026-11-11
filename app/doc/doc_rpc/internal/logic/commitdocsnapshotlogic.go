package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm/clause"
)

type CommitDocSnapshotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCommitDocSnapshotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommitDocSnapshotLogic {
	return &CommitDocSnapshotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("CommitDocSnapshotLogic", ctx),
	}
}

func (l *CommitDocSnapshotLogic) CommitDocSnapshot(in *doc_rpc.CommitDocSnapshotReq) (*doc_rpc.CommitDocSnapshotRes, error) {
	// 内部可信调用（doc_api 协同网关压实）：写全量快照 + 解码正文落权威文档 + 生成版本。
	doc, err := loadDoc(l.svcCtx.DB, in.DocId)
	if err != nil {
		return nil, errNotFound
	}
	space, err := loadSpace(l.svcCtx.DB, doc.SpaceID)
	if err != nil {
		return nil, errNotFound
	}
	epoch := doc.Epoch
	if in.Epoch > 0 {
		// 旧纪元的压实结果不能覆盖新纪元状态。
		if int(in.Epoch) != epoch {
			return nil, errInvalidArg
		}
	}

	var maxSeq int64
	l.svcCtx.DB.Model(&doc_models.DocYjsUpdate{}).
		Where("doc_id = ? AND epoch = ?", doc.DocID, epoch).
		Select("COALESCE(MAX(seq), 0)").Scan(&maxSeq)

	// 快照按 (doc_id, epoch) 唯一：不同纪元互不覆盖。
	snap := doc_models.DocYjsSnapshot{
		DocID:         doc.DocID,
		Epoch:         epoch,
		State:         in.Snapshot,
		UpdateSeqUpTo: maxSeq,
		Version:       doc.Version,
	}
	if err := l.svcCtx.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "doc_id"}, {Name: "epoch"}},
		DoUpdates: clause.AssignmentColumns([]string{"state", "update_seq_up_to", "version"}),
	}).Create(&snap).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "写入yjs快照失败", Data: err})
		return nil, err
	}

	// content_json 为空表示纯快照压实：仅收敛增量日志/写快照，不改动权威正文与版本。
	// 携带解码正文时才走 applyDocContent（更新 content + 生成版本）。
	if trim(in.ContentJson) == "" {
		return &doc_rpc.CommitDocSnapshotRes{Version: int32(doc.Version)}, nil
	}

	version, err := applyDocContent(l.ctx, l.svcCtx, doc, space, in.ContentJson, in.ContentText, in.OperatorUserId, "content.autosave")
	if err != nil {
		return nil, err
	}
	return &doc_rpc.CommitDocSnapshotRes{Version: int32(version)}, nil
}
