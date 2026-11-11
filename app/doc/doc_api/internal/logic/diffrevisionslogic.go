package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type DiffRevisionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiffRevisionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiffRevisionsLogic {
	return &DiffRevisionsLogic{ctx: ctx, svcCtx: svcCtx}
}

// DiffRevisions 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *DiffRevisionsLogic) DiffRevisions(req *types.DiffRevisionsReq) (*types.DiffRevisionsRes, error) {
	return forward[types.DiffRevisionsReq, types.DiffRevisionsRes, doc_rpc.DiffRevisionsReq, doc_rpc.DiffRevisionsRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.DiffRevisionsReq) (*doc_rpc.DiffRevisionsRes, error) {
			return l.svcCtx.Doc.DiffRevisions(ctx, in)
		},
	)
}
