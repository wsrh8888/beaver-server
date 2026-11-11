package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type DeleteDocBlockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDocBlockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDocBlockLogic {
	return &DeleteDocBlockLogic{ctx: ctx, svcCtx: svcCtx}
}

// DeleteDocBlock 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *DeleteDocBlockLogic) DeleteDocBlock(req *types.DeleteDocBlockReq) (*types.DeleteDocBlockRes, error) {
	return forward[types.DeleteDocBlockReq, types.DeleteDocBlockRes, doc_rpc.DeleteDocBlockReq, doc_rpc.DeleteDocBlockRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.DeleteDocBlockReq) (*doc_rpc.DeleteDocBlockRes, error) {
			return l.svcCtx.Doc.DeleteDocBlock(ctx, in)
		},
	)
}
