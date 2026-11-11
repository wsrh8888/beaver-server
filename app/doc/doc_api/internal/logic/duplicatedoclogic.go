package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type DuplicateDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDuplicateDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DuplicateDocLogic {
	return &DuplicateDocLogic{ctx: ctx, svcCtx: svcCtx}
}

// DuplicateDoc 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *DuplicateDocLogic) DuplicateDoc(req *types.DuplicateDocReq) (*types.DuplicateDocRes, error) {
	return forward[types.DuplicateDocReq, types.DuplicateDocRes, doc_rpc.DuplicateDocReq, doc_rpc.DuplicateDocRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.DuplicateDocReq) (*doc_rpc.DuplicateDocRes, error) {
			return l.svcCtx.Doc.DuplicateDoc(ctx, in)
		},
	)
}
