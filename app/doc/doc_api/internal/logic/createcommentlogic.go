package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{ctx: ctx, svcCtx: svcCtx}
}

// CreateComment 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *CreateCommentLogic) CreateComment(req *types.CreateCommentReq) (*types.CreateCommentRes, error) {
	return forward[types.CreateCommentReq, types.CreateCommentRes, doc_rpc.CreateCommentReq, doc_rpc.CreateCommentRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.CreateCommentReq) (*doc_rpc.CreateCommentRes, error) {
			return l.svcCtx.Doc.CreateComment(ctx, in)
		},
	)
}
