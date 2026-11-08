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

	var maxSeq int64
	l.svcCtx.DB.Model(&doc_models.DocYjsUpdate{}).Where("doc_id = ?", doc.DocID).Select("COALESCE(MAX(seq), 0)").Scan(&maxSeq)

	var snap doc_models.DocYjsSnapshot
	err = l.svcCtx.DB.Where("doc_id = ?", doc.DocID).First(&snap).Error
	if err == gorm.ErrRecordNotFound {
		snap = doc_models.DocYjsSnapshot{DocID: doc.DocID}
		err = nil
	}
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "读取yjs快照失败", Data: err})
		return nil, err
	}
	snap.State = in.Snapshot
	snap.UpdateSeqUpTo = maxSeq
	if snap.Id == 0 {
		if cerr := l.svcCtx.DB.Create(&snap).Error; cerr != nil {
			l.logger.Error(model.LogMsg{Text: "写入yjs快照失败", Data: cerr})
			return nil, cerr
		}
	} else {
		if uerr := l.svcCtx.DB.Model(&doc_models.DocYjsSnapshot{}).Where("doc_id = ?", doc.DocID).Updates(map[string]any{
			"state":            snap.State,
			"update_seq_up_to": snap.UpdateSeqUpTo,
		}).Error; uerr != nil {
			l.logger.Error(model.LogMsg{Text: "更新yjs快照失败", Data: uerr})
			return nil, uerr
		}
	}

	// content_json 为空表示纯快照压实：仅收敛增量日志/写快照，不改动权威正文与版本。
	// 携带解码正文时才走 applyDocContent（更新 content + 生成版本）。
	if trim(in.ContentJson) == "" {
		return &doc_rpc.CommitDocSnapshotRes{Version: int32(doc.Version)}, nil
	}

	version, err := applyDocContent(l.ctx, l.svcCtx, doc, space, in.ContentJson, in.ContentText, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	return &doc_rpc.CommitDocSnapshotRes{Version: int32(version)}, nil
}
