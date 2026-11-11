package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListCollaboratorsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCollaboratorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCollaboratorsLogic {
	return &ListCollaboratorsLogic{ctx: ctx, svcCtx: svcCtx}
}

// ListCollaborators 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ListCollaboratorsLogic) ListCollaborators(req *types.ListCollaboratorsReq) (*types.ListCollaboratorsRes, error) {
	return forward[types.ListCollaboratorsReq, types.ListCollaboratorsRes, doc_rpc.ListCollaboratorsReq, doc_rpc.ListCollaboratorsRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ListCollaboratorsReq) (*doc_rpc.ListCollaboratorsRes, error) {
			return l.svcCtx.Doc.ListCollaborators(ctx, in)
		},
	)
}
