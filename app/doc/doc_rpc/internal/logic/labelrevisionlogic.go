package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type LabelRevisionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewLabelRevisionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LabelRevisionLogic {
	return &LabelRevisionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("LabelRevisionLogic", ctx),
	}
}

func (l *LabelRevisionLogic) LabelRevision(in *doc_rpc.LabelRevisionReq) (*doc_rpc.LabelRevisionRes, error) {
	var rev doc_models.DocRevision
	if err := l.svcCtx.DB.Where("revision_id = ?", in.RevisionId).First(&rev).Error; err != nil {
		return nil, errNotFound
	}
	if _, _, err := requireDoc(l.svcCtx.DB, rev.DocID, in.OperatorUserId, roleEditor); err != nil {
		return nil, err
	}
	label := firstRunes(trim(in.Label), 128)
	if err := l.svcCtx.DB.Model(&doc_models.DocRevision{}).
		Where("revision_id = ?", rev.RevisionID).Update("label", label).Error; err != nil {
		return nil, err
	}
	auditJSON(l.svcCtx.DB, in.OperatorUserId, "revision.label", "doc", rev.DocID,
		map[string]any{"revisionId": rev.RevisionID, "label": label})
	return &doc_rpc.LabelRevisionRes{}, nil
}
