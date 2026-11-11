package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListDocsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDocsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDocsLogic {
	return &ListDocsLogic{ctx: ctx, svcCtx: svcCtx}
}

// ListDocs 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ListDocsLogic) ListDocs(req *types.ListDocsReq) (*types.ListDocsRes, error) {
	return forward[types.ListDocsReq, types.ListDocsRes, doc_rpc.ListDocsReq, doc_rpc.ListDocsRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ListDocsReq) (*doc_rpc.ListDocsRes, error) {
			return l.svcCtx.Doc.ListDocs(ctx, in)
		},
	)
}
