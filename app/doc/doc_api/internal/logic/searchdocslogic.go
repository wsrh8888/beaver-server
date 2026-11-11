package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type SearchDocsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchDocsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchDocsLogic {
	return &SearchDocsLogic{ctx: ctx, svcCtx: svcCtx}
}

// SearchDocs 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *SearchDocsLogic) SearchDocs(req *types.SearchDocsReq) (*types.SearchDocsRes, error) {
	return forward[types.SearchDocsReq, types.SearchDocsRes, doc_rpc.SearchDocsReq, doc_rpc.SearchDocsRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.SearchDocsReq) (*doc_rpc.SearchDocsRes, error) {
			return l.svcCtx.Doc.SearchDocs(ctx, in)
		},
	)
}
