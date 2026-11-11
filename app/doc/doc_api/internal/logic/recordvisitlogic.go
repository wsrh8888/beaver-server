package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type RecordVisitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecordVisitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordVisitLogic {
	return &RecordVisitLogic{ctx: ctx, svcCtx: svcCtx}
}

// RecordVisit 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *RecordVisitLogic) RecordVisit(req *types.RecordVisitReq) (*types.RecordVisitRes, error) {
	return forward[types.RecordVisitReq, types.RecordVisitRes, doc_rpc.RecordVisitReq, doc_rpc.RecordVisitRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.RecordVisitReq) (*doc_rpc.RecordVisitRes, error) {
			return l.svcCtx.Doc.RecordVisit(ctx, in)
		},
	)
}
