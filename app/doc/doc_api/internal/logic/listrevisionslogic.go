package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListRevisionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRevisionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRevisionsLogic {
	return &ListRevisionsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListRevisionsLogic) ListRevisions(req *types.ListRevisionsReq) (resp *types.ListRevisionsRes, err error) {
	// 版本历史可读需 viewer 权限（RPC 无主体，API 层校验）。
	if err := requireRole(l.ctx, l.svcCtx, "doc", req.DocId, req.UserId, "viewer"); err != nil {
		return nil, err
	}
	res, err := l.svcCtx.Doc.ListRevisions(l.ctx, &doc_rpc.ListRevisionsReq{
		DocId:    req.DocId,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	return &types.ListRevisionsRes{Total: res.Total, Revisions: convRevisionItems(res.Revisions)}, nil
}
