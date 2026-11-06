package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type UpdateNodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNodeLogic {
	return &UpdateNodeLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateNodeLogic) UpdateNode(req *types.UpdateNodeReq) (resp *types.UpdateNodeRes, err error) {
	_, err = l.svcCtx.Doc.UpdateNode(l.ctx, &doc_rpc.UpdateNodeReq{
		NodeId:         req.NodeId,
		Title:          req.Title,
		Icon:           req.Icon,
		ParentId:       req.ParentId,
		Sort:           int32(req.Sort),
		OperatorUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.UpdateNodeRes{}, nil
}
