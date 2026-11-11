package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type SaveDocBlockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveDocBlockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDocBlockLogic {
	return &SaveDocBlockLogic{ctx: ctx, svcCtx: svcCtx}
}

// SaveDocBlock 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *SaveDocBlockLogic) SaveDocBlock(req *types.SaveDocBlockReq) (*types.SaveDocBlockRes, error) {
	return forward[types.SaveDocBlockReq, types.SaveDocBlockRes, doc_rpc.SaveDocBlockReq, doc_rpc.SaveDocBlockRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.SaveDocBlockReq) (*doc_rpc.SaveDocBlockRes, error) {
			return l.svcCtx.Doc.SaveDocBlock(ctx, in)
		},
	)
}
