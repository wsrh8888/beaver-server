package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type EmptyTrashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmptyTrashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmptyTrashLogic {
	return &EmptyTrashLogic{ctx: ctx, svcCtx: svcCtx}
}

// EmptyTrash 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *EmptyTrashLogic) EmptyTrash(req *types.EmptyTrashReq) (*types.EmptyTrashRes, error) {
	return forward[types.EmptyTrashReq, types.EmptyTrashRes, doc_rpc.EmptyTrashReq, doc_rpc.EmptyTrashRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.EmptyTrashReq) (*doc_rpc.EmptyTrashRes, error) {
			return l.svcCtx.Doc.EmptyTrash(ctx, in)
		},
	)
}
