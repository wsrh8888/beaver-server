package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type CreateSpaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSpaceLogic {
	return &CreateSpaceLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateSpaceLogic) CreateSpace(req *types.CreateSpaceReq) (resp *types.CreateSpaceRes, err error) {
	res, err := l.svcCtx.Doc.CreateSpace(l.ctx, &doc_rpc.CreateSpaceReq{
		Name:        req.Name,
		Kind:        req.Kind,
		OwnerUserId: req.UserId,
		ParentId:    req.ParentId,
		Icon:        req.Icon,
		Description: req.Description,
		Sort:        int32(req.Sort),
	})
	if err != nil {
		return nil, err
	}
	return &types.CreateSpaceRes{SpaceId: res.SpaceId}, nil
}
