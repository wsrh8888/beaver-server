package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type CreateNodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNodeLogic {
	return &CreateNodeLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateNodeLogic) CreateNode(req *types.CreateNodeReq) (resp *types.CreateNodeRes, err error) {
	res, err := l.svcCtx.Doc.CreateNode(l.ctx, &doc_rpc.CreateNodeReq{
		SpaceId:        req.SpaceId,
		ParentId:       req.ParentId,
		Title:          req.Title,
		Type:           req.Type,
		Icon:           req.Icon,
		Sort:           int32(req.Sort),
		OperatorUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.CreateNodeRes{NodeId: res.NodeId, DocId: res.DocId}, nil
}
