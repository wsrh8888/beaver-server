package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListAuditLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAuditLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAuditLogsLogic {
	return &ListAuditLogsLogic{ctx: ctx, svcCtx: svcCtx}
}

// ListAuditLogs 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ListAuditLogsLogic) ListAuditLogs(req *types.ListAuditLogsReq) (*types.ListAuditLogsRes, error) {
	return forward[types.ListAuditLogsReq, types.ListAuditLogsRes, doc_rpc.ListAuditLogsReq, doc_rpc.ListAuditLogsRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ListAuditLogsReq) (*doc_rpc.ListAuditLogsRes, error) {
			return l.svcCtx.Doc.ListAuditLogs(ctx, in)
		},
	)
}
