package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type UpdateSpaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSpaceLogic {
	return &UpdateSpaceLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateSpaceLogic) UpdateSpace(req *types.UpdateSpaceReq) (resp *types.UpdateSpaceRes, err error) {
	// UpdateSpace RPC 报文不含 operator，权限在 API 层补齐：需空间管理权。
	if err := requireRole(l.ctx, l.svcCtx, "space", req.SpaceId, req.UserId, "manager"); err != nil {
		return nil, err
	}
	_, err = l.svcCtx.Doc.UpdateSpace(l.ctx, &doc_rpc.UpdateSpaceReq{
		SpaceId:     req.SpaceId,
		Name:        req.Name,
		Icon:        req.Icon,
		Description: req.Description,
		ParentId:    req.ParentId,
		Status:      int32(req.Status),
		Sort:        int32(req.Sort),
	})
	if err != nil {
		return nil, err
	}
	return &types.UpdateSpaceRes{}, nil
}
