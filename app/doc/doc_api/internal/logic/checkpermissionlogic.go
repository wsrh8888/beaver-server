package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type CheckPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPermissionLogic {
	return &CheckPermissionLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CheckPermissionLogic) CheckPermission(req *types.CheckPermissionReq) (resp *types.CheckPermissionRes, err error) {
	res, err := l.svcCtx.Doc.CheckPermission(l.ctx, &doc_rpc.CheckPermissionReq{
		ResourceType: req.ResourceType,
		ResourceId:   req.ResourceId,
		UserId:       req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.CheckPermissionRes{Role: res.Role}, nil
}
