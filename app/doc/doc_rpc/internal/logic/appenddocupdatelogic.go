package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type AppendDocUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewAppendDocUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppendDocUpdateLogic {
	return &AppendDocUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("AppendDocUpdateLogic", ctx),
	}
}

func (l *AppendDocUpdateLogic) AppendDocUpdate(in *doc_rpc.AppendDocUpdateReq) (*doc_rpc.AppendDocUpdateRes, error) {
	// 内部可信调用（doc_sync_ws）。文档内 seq 递增，作为增量日志重放顺序。
	var maxSeq int64
	l.svcCtx.DB.Model(&doc_models.DocYjsUpdate{}).Where("doc_id = ?", in.DocId).Select("COALESCE(MAX(seq), 0)").Scan(&maxSeq)

	row := doc_models.DocYjsUpdate{
		DocID:        in.DocId,
		Seq:          maxSeq + 1,
		Update:       in.Update,
		OriginUserID: trim(in.OriginUserId),
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "追加yjs增量失败", Data: err})
		return nil, err
	}
	return &doc_rpc.AppendDocUpdateRes{Seq: row.Seq}, nil
}
