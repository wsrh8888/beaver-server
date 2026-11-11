package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ExportDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportDocLogic {
	return &ExportDocLogic{ctx: ctx, svcCtx: svcCtx}
}

// ExportDoc 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ExportDocLogic) ExportDoc(req *types.ExportDocReq) (*types.ExportDocRes, error) {
	return forward[types.ExportDocReq, types.ExportDocRes, doc_rpc.ExportDocReq, doc_rpc.ExportDocRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ExportDocReq) (*doc_rpc.ExportDocRes, error) {
			return l.svcCtx.Doc.ExportDoc(ctx, in)
		},
	)
}
