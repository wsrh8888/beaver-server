package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListDocBlocksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDocBlocksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDocBlocksLogic {
	return &ListDocBlocksLogic{ctx: ctx, svcCtx: svcCtx}
}

// ListDocBlocks 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ListDocBlocksLogic) ListDocBlocks(req *types.ListDocBlocksReq) (*types.ListDocBlocksRes, error) {
	return forward[types.ListDocBlocksReq, types.ListDocBlocksRes, doc_rpc.ListDocBlocksReq, doc_rpc.ListDocBlocksRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ListDocBlocksReq) (*doc_rpc.ListDocBlocksRes, error) {
			return l.svcCtx.Doc.ListDocBlocks(ctx, in)
		},
	)
}
