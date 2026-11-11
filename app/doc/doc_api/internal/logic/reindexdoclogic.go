package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ReindexDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReindexDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReindexDocLogic {
	return &ReindexDocLogic{ctx: ctx, svcCtx: svcCtx}
}

// ReindexDoc 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ReindexDocLogic) ReindexDoc(req *types.ReindexDocReq) (*types.ReindexDocRes, error) {
	return forward[types.ReindexDocReq, types.ReindexDocRes, doc_rpc.ReindexDocReq, doc_rpc.ReindexDocRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ReindexDocReq) (*doc_rpc.ReindexDocRes, error) {
			return l.svcCtx.Doc.ReindexDoc(ctx, in)
		},
	)
}
