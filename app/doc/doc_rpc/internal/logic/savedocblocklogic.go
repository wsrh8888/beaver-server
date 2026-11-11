package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type SaveDocBlockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSaveDocBlockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDocBlockLogic {
	return &SaveDocBlockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("SaveDocBlockLogic", ctx),
	}
}

func (l *SaveDocBlockLogic) SaveDocBlock(in *doc_rpc.SaveDocBlockReq) (*doc_rpc.SaveDocBlockRes, error) {
	doc, _, err := requireDoc(l.svcCtx.DB, in.DocId, in.OperatorUserId, roleEditor)
	if err != nil {
		return nil, err
	}
	kind := trim(in.Kind)
	if kind != "sheet" && kind != "board" {
		return nil, errInvalidArg
	}
	actor := trim(in.OperatorUserId)
	var row doc_models.DocBlock
	if bid := trim(in.BlockId); bid != "" {
		if err := l.svcCtx.DB.Where("block_id = ? AND doc_id = ?", bid, doc.DocID).
			First(&row).Error; err != nil {
			return nil, errNotFound
		}
		row.Version++
	} else {
		row = doc_models.DocBlock{BlockID: genID("blk"), DocID: doc.DocID, Version: 1}
	}
	row.Kind = kind
	row.Title = trim(in.Title)
	row.PayloadJSON = in.PayloadJson
	row.UpdatedBy = actor
	if err := l.svcCtx.DB.Save(&row).Error; err != nil {
		return nil, err
	}
	auditJSON(l.svcCtx.DB, actor, "block.save", "doc", doc.DocID,
		map[string]any{"blockId": row.BlockID, "kind": kind})
	return &doc_rpc.SaveDocBlockRes{Block: toDocBlockItem(&row)}, nil
}
