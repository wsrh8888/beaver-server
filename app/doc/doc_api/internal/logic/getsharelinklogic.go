package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type GetShareLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShareLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShareLinkLogic {
	return &GetShareLinkLogic{ctx: ctx, svcCtx: svcCtx}
}

// GetShareLink 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *GetShareLinkLogic) GetShareLink(req *types.GetShareLinkReq) (*types.GetShareLinkRes, error) {
	return forward[types.GetShareLinkReq, types.GetShareLinkRes, doc_rpc.GetShareLinkReq, doc_rpc.GetShareLinkRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.GetShareLinkReq) (*doc_rpc.GetShareLinkRes, error) {
			return l.svcCtx.Doc.GetShareLink(ctx, in)
		},
	)
}
