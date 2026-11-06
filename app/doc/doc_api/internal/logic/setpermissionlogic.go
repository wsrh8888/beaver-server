package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type SetPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetPermissionLogic {
	return &SetPermissionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *SetPermissionLogic) SetPermission(req *types.SetPermissionReq) (resp *types.SetPermissionRes, err error) {
	_, err = l.svcCtx.Doc.SetPermission(l.ctx, &doc_rpc.SetPermissionReq{
		ResourceType:  req.ResourceType,
		ResourceId:    req.ResourceId,
		PrincipalType: req.PrincipalType,
		PrincipalId:   req.PrincipalId,
		Role:          req.Role,
		OperatorUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.SetPermissionRes{}, nil
}
