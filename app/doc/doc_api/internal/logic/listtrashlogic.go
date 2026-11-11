package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListTrashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTrashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTrashLogic {
	return &ListTrashLogic{ctx: ctx, svcCtx: svcCtx}
}

// ListTrash 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ListTrashLogic) ListTrash(req *types.ListTrashReq) (*types.ListTrashRes, error) {
	return forward[types.ListTrashReq, types.ListTrashRes, doc_rpc.ListTrashReq, doc_rpc.ListTrashRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ListTrashReq) (*doc_rpc.ListTrashRes, error) {
			return l.svcCtx.Doc.ListTrash(ctx, in)
		},
	)
}
