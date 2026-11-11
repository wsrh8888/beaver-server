package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type SetShareLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetShareLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetShareLinkLogic {
	return &SetShareLinkLogic{ctx: ctx, svcCtx: svcCtx}
}

// SetShareLink 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *SetShareLinkLogic) SetShareLink(req *types.SetShareLinkReq) (*types.SetShareLinkRes, error) {
	return forward[types.SetShareLinkReq, types.SetShareLinkRes, doc_rpc.SetShareLinkReq, doc_rpc.SetShareLinkRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.SetShareLinkReq) (*doc_rpc.SetShareLinkRes, error) {
			return l.svcCtx.Doc.SetShareLink(ctx, in)
		},
	)
}
