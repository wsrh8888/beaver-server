package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ResolveShareLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResolveShareLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResolveShareLinkLogic {
	return &ResolveShareLinkLogic{ctx: ctx, svcCtx: svcCtx}
}

// ResolveShareLink 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ResolveShareLinkLogic) ResolveShareLink(req *types.ResolveShareLinkReq) (*types.ResolveShareLinkRes, error) {
	return forward[types.ResolveShareLinkReq, types.ResolveShareLinkRes, doc_rpc.ResolveShareLinkReq, doc_rpc.ResolveShareLinkRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ResolveShareLinkReq) (*doc_rpc.ResolveShareLinkRes, error) {
			return l.svcCtx.Doc.ResolveShareLink(ctx, in)
		},
	)
}
