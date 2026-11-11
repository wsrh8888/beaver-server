package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type RevokeShareLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRevokeShareLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeShareLinkLogic {
	return &RevokeShareLinkLogic{ctx: ctx, svcCtx: svcCtx}
}

// RevokeShareLink 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *RevokeShareLinkLogic) RevokeShareLink(req *types.RevokeShareLinkReq) (*types.RevokeShareLinkRes, error) {
	return forward[types.RevokeShareLinkReq, types.RevokeShareLinkRes, doc_rpc.RevokeShareLinkReq, doc_rpc.RevokeShareLinkRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.RevokeShareLinkReq) (*doc_rpc.RevokeShareLinkRes, error) {
			return l.svcCtx.Doc.RevokeShareLink(ctx, in)
		},
	)
}
