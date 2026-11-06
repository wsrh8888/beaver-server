package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type GetRevisionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRevisionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRevisionLogic {
	return &GetRevisionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetRevisionLogic) GetRevision(req *types.GetRevisionReq) (resp *types.GetRevisionRes, err error) {
	res, err := l.svcCtx.Doc.GetRevision(l.ctx, &doc_rpc.GetRevisionReq{
		RevisionId: req.RevisionId,
	})
	if err != nil {
		return nil, err
	}
	if res.Revision == nil {
		return nil, errNoPermission
	}
	// 取到版本归属文档后校验读权限，未授权不外泄 content_json。
	if err := requireRole(l.ctx, l.svcCtx, "doc", res.Revision.DocId, req.UserId, "viewer"); err != nil {
		return nil, err
	}
	return &types.GetRevisionRes{
		Revision:    revisionPtr(res.Revision),
		ContentJson: res.ContentJson,
	}, nil
}

func revisionPtr(it *doc_rpc.RevisionItem) *types.RevisionItem {
	if it == nil {
		return nil
	}
	v := convRevisionItem(it)
	return &v
}
