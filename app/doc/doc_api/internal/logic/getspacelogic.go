package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type GetSpaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSpaceLogic {
	return &GetSpaceLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetSpaceLogic) GetSpace(req *types.GetSpaceReq) (resp *types.GetSpaceRes, err error) {
	res, err := l.svcCtx.Doc.GetSpace(l.ctx, &doc_rpc.GetSpaceReq{
		SpaceId: req.SpaceId,
		UserId:  req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.GetSpaceRes{Space: spacePtr(res.Space)}, nil
}

func spacePtr(it *doc_rpc.SpaceItem) *types.SpaceItem {
	if it == nil {
		return nil
	}
	v := convSpaceItem(it)
	return &v
}
