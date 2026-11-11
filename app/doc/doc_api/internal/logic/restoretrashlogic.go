package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type RestoreTrashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRestoreTrashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreTrashLogic {
	return &RestoreTrashLogic{ctx: ctx, svcCtx: svcCtx}
}

// RestoreTrash 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *RestoreTrashLogic) RestoreTrash(req *types.RestoreTrashReq) (*types.RestoreTrashRes, error) {
	return forward[types.RestoreTrashReq, types.RestoreTrashRes, doc_rpc.RestoreTrashReq, doc_rpc.RestoreTrashRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.RestoreTrashReq) (*doc_rpc.RestoreTrashRes, error) {
			return l.svcCtx.Doc.RestoreTrash(ctx, in)
		},
	)
}
