package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListTemplatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTemplatesLogic {
	return &ListTemplatesLogic{ctx: ctx, svcCtx: svcCtx}
}

// ListTemplates 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ListTemplatesLogic) ListTemplates(req *types.ListTemplatesReq) (*types.ListTemplatesRes, error) {
	return forward[types.ListTemplatesReq, types.ListTemplatesRes, doc_rpc.ListTemplatesReq, doc_rpc.ListTemplatesRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ListTemplatesReq) (*doc_rpc.ListTemplatesRes, error) {
			return l.svcCtx.Doc.ListTemplates(ctx, in)
		},
	)
}
