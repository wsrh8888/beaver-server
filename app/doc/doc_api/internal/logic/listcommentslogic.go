package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{ctx: ctx, svcCtx: svcCtx}
}

// ListComments 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ListCommentsLogic) ListComments(req *types.ListCommentsReq) (*types.ListCommentsRes, error) {
	return forward[types.ListCommentsReq, types.ListCommentsRes, doc_rpc.ListCommentsReq, doc_rpc.ListCommentsRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ListCommentsReq) (*doc_rpc.ListCommentsRes, error) {
			return l.svcCtx.Doc.ListComments(ctx, in)
		},
	)
}
