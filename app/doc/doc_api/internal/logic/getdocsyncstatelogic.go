package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type GetDocSyncStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDocSyncStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDocSyncStateLogic {
	return &GetDocSyncStateLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetDocSyncStateLogic) GetDocSyncState(req *types.GetDocSyncStateReq) (resp *types.GetDocSyncStateRes, err error) {
	// doc_rpc 的同步原语面向可信网关，用户面 HTTP 入口在此补齐读权限校验。
	if err := requireRole(l.ctx, l.svcCtx, "doc", req.DocId, req.UserId, "viewer"); err != nil {
		return nil, err
	}
	res, err := l.svcCtx.Doc.GetDocSyncState(l.ctx, &doc_rpc.GetDocSyncStateReq{
		DocId: req.DocId,
	})
	if err != nil {
		return nil, err
	}
	return &types.GetDocSyncStateRes{
		Snapshot: res.Snapshot,
		Updates:  res.Updates,
		Version:  int(res.Version),
	}, nil
}
