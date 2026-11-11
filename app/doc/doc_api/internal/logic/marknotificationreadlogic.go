package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type MarkNotificationReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarkNotificationReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkNotificationReadLogic {
	return &MarkNotificationReadLogic{ctx: ctx, svcCtx: svcCtx}
}

// MarkNotificationRead 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *MarkNotificationReadLogic) MarkNotificationRead(req *types.MarkNotificationReadReq) (*types.MarkNotificationReadRes, error) {
	return forward[types.MarkNotificationReadReq, types.MarkNotificationReadRes, doc_rpc.MarkNotificationReadReq, doc_rpc.MarkNotificationReadRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.MarkNotificationReadReq) (*doc_rpc.MarkNotificationReadRes, error) {
			return l.svcCtx.Doc.MarkNotificationRead(ctx, in)
		},
	)
}
