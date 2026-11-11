package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type DeleteDocBlockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteDocBlockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDocBlockLogic {
	return &DeleteDocBlockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DeleteDocBlockLogic", ctx),
	}
}

func (l *DeleteDocBlockLogic) DeleteDocBlock(in *doc_rpc.DeleteDocBlockReq) (*doc_rpc.DeleteDocBlockRes, error) {
	doc, _, err := requireDoc(l.svcCtx.DB, in.DocId, in.OperatorUserId, roleEditor)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.DB.Where("doc_id = ? AND block_id = ?", doc.DocID, in.BlockId).
		Delete(&doc_models.DocBlock{}).Error; err != nil {
		return nil, err
	}
	auditJSON(l.svcCtx.DB, in.OperatorUserId, "block.delete", "doc", doc.DocID,
		map[string]any{"blockId": in.BlockId})
	return &doc_rpc.DeleteDocBlockRes{}, nil
}
