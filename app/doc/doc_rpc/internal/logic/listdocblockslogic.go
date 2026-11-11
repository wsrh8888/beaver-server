package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ListDocBlocksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListDocBlocksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDocBlocksLogic {
	return &ListDocBlocksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListDocBlocksLogic", ctx),
	}
}

func (l *ListDocBlocksLogic) ListDocBlocks(in *doc_rpc.ListDocBlocksReq) (*doc_rpc.ListDocBlocksRes, error) {
	doc, _, err := requireDoc(l.svcCtx.DB, in.DocId, in.UserId, roleViewer)
	if err != nil {
		return nil, err
	}
	var rows []doc_models.DocBlock
	if err := l.svcCtx.DB.Where("doc_id = ?", doc.DocID).Order("id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	blocks := make([]*doc_rpc.DocBlockItem, 0, len(rows))
	for i := range rows {
		blocks = append(blocks, toDocBlockItem(&rows[i]))
	}
	return &doc_rpc.ListDocBlocksRes{Blocks: blocks}, nil
}
