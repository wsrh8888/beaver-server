package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListNotificationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotificationsLogic {
	return &ListNotificationsLogic{ctx: ctx, svcCtx: svcCtx}
}

// ListNotifications 把 HTTP 请求转成 doc_rpc 调用。字段映射与鉴权主体注入由 forward 统一处理。
func (l *ListNotificationsLogic) ListNotifications(req *types.ListNotificationsReq) (*types.ListNotificationsRes, error) {
	return forward[types.ListNotificationsReq, types.ListNotificationsRes, doc_rpc.ListNotificationsReq, doc_rpc.ListNotificationsRes](
		l.ctx, req,
		func(ctx context.Context, in *doc_rpc.ListNotificationsReq) (*doc_rpc.ListNotificationsRes, error) {
			return l.svcCtx.Doc.ListNotifications(ctx, in)
		},
	)
}
