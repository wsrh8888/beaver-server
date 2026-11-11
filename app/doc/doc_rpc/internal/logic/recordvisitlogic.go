package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"

	"gorm.io/gorm/clause"
)

type RecordVisitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRecordVisitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordVisitLogic {
	return &RecordVisitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("RecordVisitLogic", ctx),
	}
}

func (l *RecordVisitLogic) RecordVisit(in *doc_rpc.RecordVisitReq) (*doc_rpc.RecordVisitRes, error) {
	doc, _, err := requireDoc(l.svcCtx.DB, in.DocId, in.UserId, roleViewer)
	if err != nil {
		return nil, err
	}
	row := doc_models.DocRecentVisit{
		UserID:    trim(in.UserId),
		DocID:     doc.DocID,
		VisitedAt: nowStr(),
	}
	if err := l.svcCtx.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "doc_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"visited_at"}),
	}).Create(&row).Error; err != nil {
		return nil, err
	}
	return &doc_rpc.RecordVisitRes{}, nil
}
