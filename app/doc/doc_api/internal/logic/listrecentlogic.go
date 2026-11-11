package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListRecentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRecentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRecentLogic {
	return &ListRecentLogic{ctx: ctx, svcCtx: svcCtx}
}

// ListRecent 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ListRecentLogic) ListRecent(req *types.ListRecentReq) (*types.ListRecentRes, error) {
	return forward[types.ListRecentReq, types.ListRecentRes, doc_rpc.ListRecentReq, doc_rpc.ListRecentRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ListRecentReq) (*doc_rpc.ListRecentRes, error) {
			return l.svcCtx.Doc.ListRecent(ctx, in)
		},
	)
}
