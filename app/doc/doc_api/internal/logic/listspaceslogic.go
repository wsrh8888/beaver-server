package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListSpacesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSpacesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSpacesLogic {
	return &ListSpacesLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListSpacesLogic) ListSpaces(req *types.ListSpacesReq) (resp *types.ListSpacesRes, err error) {
	res, err := l.svcCtx.Doc.ListSpaces(l.ctx, &doc_rpc.ListSpacesReq{
		UserId:   req.UserId,
		Status:   int32(req.Status),
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	return &types.ListSpacesRes{Total: res.Total, Spaces: convSpaceItems(res.Spaces)}, nil
}
