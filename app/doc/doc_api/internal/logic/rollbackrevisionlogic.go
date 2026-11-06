package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type RollbackRevisionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRollbackRevisionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RollbackRevisionLogic {
	return &RollbackRevisionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *RollbackRevisionLogic) RollbackRevision(req *types.RollbackRevisionReq) (resp *types.RollbackRevisionRes, err error) {
	res, err := l.svcCtx.Doc.RollbackRevision(l.ctx, &doc_rpc.RollbackRevisionReq{
		DocId:          req.DocId,
		RevisionId:     req.RevisionId,
		OperatorUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.RollbackRevisionRes{Version: int(res.Version)}, nil
}
