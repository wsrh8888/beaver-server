package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{ctx: ctx, svcCtx: svcCtx}
}

// DeleteComment 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *DeleteCommentLogic) DeleteComment(req *types.DeleteCommentReq) (*types.DeleteCommentRes, error) {
	return forward[types.DeleteCommentReq, types.DeleteCommentRes, doc_rpc.DeleteCommentReq, doc_rpc.DeleteCommentRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.DeleteCommentReq) (*doc_rpc.DeleteCommentRes, error) {
			return l.svcCtx.Doc.DeleteComment(ctx, in)
		},
	)
}
