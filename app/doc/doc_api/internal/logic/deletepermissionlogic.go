package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type DeletePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePermissionLogic {
	return &DeletePermissionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeletePermissionLogic) DeletePermission(req *types.DeletePermissionReq) (resp *types.DeletePermissionRes, err error) {
	_, err = l.svcCtx.Doc.DeletePermission(l.ctx, &doc_rpc.DeletePermissionReq{
		ResourceType:   req.ResourceType,
		ResourceId:     req.ResourceId,
		PrincipalType:  req.PrincipalType,
		PrincipalId:    req.PrincipalId,
		OperatorUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.DeletePermissionRes{}, nil
}
