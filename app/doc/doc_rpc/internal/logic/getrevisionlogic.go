package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
)

type GetRevisionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetRevisionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRevisionLogic {
	return &GetRevisionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetRevisionLogic", ctx),
	}
}

func (l *GetRevisionLogic) GetRevision(in *doc_rpc.GetRevisionReq) (*doc_rpc.GetRevisionRes, error) {
	var row doc_models.DocRevision
	if err := l.svcCtx.DB.Where("revision_id = ?", in.RevisionId).First(&row).Error; err != nil {
		return nil, errNotFound
	}
	return &doc_rpc.GetRevisionRes{
		Revision:    toRevisionItem(&row),
		ContentJson: row.ContentJSON,
	}, nil
}
