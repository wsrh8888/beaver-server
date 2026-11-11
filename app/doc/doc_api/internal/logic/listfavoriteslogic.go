package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListFavoritesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFavoritesLogic {
	return &ListFavoritesLogic{ctx: ctx, svcCtx: svcCtx}
}

// ListFavorites 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ListFavoritesLogic) ListFavorites(req *types.ListFavoritesReq) (*types.ListFavoritesRes, error) {
	return forward[types.ListFavoritesReq, types.ListFavoritesRes, doc_rpc.ListFavoritesReq, doc_rpc.ListFavoritesRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ListFavoritesReq) (*doc_rpc.ListFavoritesRes, error) {
			return l.svcCtx.Doc.ListFavorites(ctx, in)
		},
	)
}
