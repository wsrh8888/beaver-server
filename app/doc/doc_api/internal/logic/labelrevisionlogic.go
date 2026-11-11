package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type LabelRevisionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLabelRevisionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LabelRevisionLogic {
	return &LabelRevisionLogic{ctx: ctx, svcCtx: svcCtx}
}

// LabelRevision 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *LabelRevisionLogic) LabelRevision(req *types.LabelRevisionReq) (*types.LabelRevisionRes, error) {
	return forward[types.LabelRevisionReq, types.LabelRevisionRes, doc_rpc.LabelRevisionReq, doc_rpc.LabelRevisionRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.LabelRevisionReq) (*doc_rpc.LabelRevisionRes, error) {
			return l.svcCtx.Doc.LabelRevision(ctx, in)
		},
	)
}
