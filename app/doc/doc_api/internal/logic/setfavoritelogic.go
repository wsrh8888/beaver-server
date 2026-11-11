package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type SetFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetFavoriteLogic {
	return &SetFavoriteLogic{ctx: ctx, svcCtx: svcCtx}
}

// SetFavorite 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *SetFavoriteLogic) SetFavorite(req *types.SetFavoriteReq) (*types.SetFavoriteRes, error) {
	return forward[types.SetFavoriteReq, types.SetFavoriteRes, doc_rpc.SetFavoriteReq, doc_rpc.SetFavoriteRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.SetFavoriteReq) (*doc_rpc.SetFavoriteRes, error) {
			return l.svcCtx.Doc.SetFavorite(ctx, in)
		},
	)
}
