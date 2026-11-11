package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ResolveCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResolveCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResolveCommentLogic {
	return &ResolveCommentLogic{ctx: ctx, svcCtx: svcCtx}
}

// ResolveComment 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ResolveCommentLogic) ResolveComment(req *types.ResolveCommentReq) (*types.ResolveCommentRes, error) {
	return forward[types.ResolveCommentReq, types.ResolveCommentRes, doc_rpc.ResolveCommentReq, doc_rpc.ResolveCommentRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ResolveCommentReq) (*doc_rpc.ResolveCommentRes, error) {
			return l.svcCtx.Doc.ResolveComment(ctx, in)
		},
	)
}
