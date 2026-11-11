package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type PurgeTrashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPurgeTrashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PurgeTrashLogic {
	return &PurgeTrashLogic{ctx: ctx, svcCtx: svcCtx}
}

// PurgeTrash 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *PurgeTrashLogic) PurgeTrash(req *types.PurgeTrashReq) (*types.PurgeTrashRes, error) {
	return forward[types.PurgeTrashReq, types.PurgeTrashRes, doc_rpc.PurgeTrashReq, doc_rpc.PurgeTrashRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.PurgeTrashReq) (*doc_rpc.PurgeTrashRes, error) {
			return l.svcCtx.Doc.PurgeTrash(ctx, in)
		},
	)
}
