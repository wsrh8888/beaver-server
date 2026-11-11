package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type InstantiateTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInstantiateTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstantiateTemplateLogic {
	return &InstantiateTemplateLogic{ctx: ctx, svcCtx: svcCtx}
}

// InstantiateTemplate 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *InstantiateTemplateLogic) InstantiateTemplate(req *types.InstantiateTemplateReq) (*types.InstantiateTemplateRes, error) {
	return forward[types.InstantiateTemplateReq, types.InstantiateTemplateRes, doc_rpc.InstantiateTemplateReq, doc_rpc.InstantiateTemplateRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.InstantiateTemplateReq) (*doc_rpc.InstantiateTemplateRes, error) {
			return l.svcCtx.Doc.InstantiateTemplate(ctx, in)
		},
	)
}
