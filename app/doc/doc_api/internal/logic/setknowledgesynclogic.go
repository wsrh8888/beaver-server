package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type SetKnowledgeSyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetKnowledgeSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetKnowledgeSyncLogic {
	return &SetKnowledgeSyncLogic{ctx: ctx, svcCtx: svcCtx}
}

// SetKnowledgeSync 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *SetKnowledgeSyncLogic) SetKnowledgeSync(req *types.SetKnowledgeSyncReq) (*types.SetKnowledgeSyncRes, error) {
	return forward[types.SetKnowledgeSyncReq, types.SetKnowledgeSyncRes, doc_rpc.SetKnowledgeSyncReq, doc_rpc.SetKnowledgeSyncRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.SetKnowledgeSyncReq) (*doc_rpc.SetKnowledgeSyncRes, error) {
			return l.svcCtx.Doc.SetKnowledgeSync(ctx, in)
		},
	)
}
