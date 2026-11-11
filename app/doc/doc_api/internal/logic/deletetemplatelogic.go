package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type DeleteTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTemplateLogic {
	return &DeleteTemplateLogic{ctx: ctx, svcCtx: svcCtx}
}

// DeleteTemplate 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *DeleteTemplateLogic) DeleteTemplate(req *types.DeleteTemplateReq) (*types.DeleteTemplateRes, error) {
	return forward[types.DeleteTemplateReq, types.DeleteTemplateRes, doc_rpc.DeleteTemplateReq, doc_rpc.DeleteTemplateRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.DeleteTemplateReq) (*doc_rpc.DeleteTemplateRes, error) {
			return l.svcCtx.Doc.DeleteTemplate(ctx, in)
		},
	)
}
