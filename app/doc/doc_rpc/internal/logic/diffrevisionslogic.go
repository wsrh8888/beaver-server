package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type DiffRevisionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDiffRevisionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiffRevisionsLogic {
	return &DiffRevisionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DiffRevisionsLogic", ctx),
	}
}

func (l *DiffRevisionsLogic) DiffRevisions(in *doc_rpc.DiffRevisionsReq) (*doc_rpc.DiffRevisionsRes, error) {
	var base, target doc_models.DocRevision
	if err := l.svcCtx.DB.Where("revision_id = ?", in.BaseRevisionId).First(&base).Error; err != nil {
		return nil, errNotFound
	}
	if err := l.svcCtx.DB.Where("revision_id = ?", in.TargetRevisionId).First(&target).Error; err != nil {
		return nil, errNotFound
	}
	if base.DocID != target.DocID {
		return nil, errInvalidArg
	}
	if _, _, err := requireDoc(l.svcCtx.DB, base.DocID, in.UserId, roleViewer); err != nil {
		return nil, err
	}
	baseText := extractTextFromJSON(base.ContentJSON)
	targetText := extractTextFromJSON(target.ContentJSON)
	added, removed := lineDiff(baseText, targetText)
	return &doc_rpc.DiffRevisionsRes{
		BaseContentJson:   base.ContentJSON,
		TargetContentJson: target.ContentJSON,
		BaseText:          baseText,
		TargetText:        targetText,
		AddedLines:        added,
		RemovedLines:      removed,
	}, nil
}
