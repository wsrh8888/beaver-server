package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type CreateTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTemplateLogic {
	return &CreateTemplateLogic{ctx: ctx, svcCtx: svcCtx}
}

// CreateTemplate 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *CreateTemplateLogic) CreateTemplate(req *types.CreateTemplateReq) (*types.CreateTemplateRes, error) {
	return forward[types.CreateTemplateReq, types.CreateTemplateRes, doc_rpc.CreateTemplateReq, doc_rpc.CreateTemplateRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.CreateTemplateReq) (*doc_rpc.CreateTemplateRes, error) {
			return l.svcCtx.Doc.CreateTemplate(ctx, in)
		},
	)
}
