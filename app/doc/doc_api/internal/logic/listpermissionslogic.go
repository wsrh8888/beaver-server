package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type ListPermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPermissionsLogic {
	return &ListPermissionsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListPermissionsLogic) ListPermissions(req *types.ListPermissionsReq) (resp *types.ListPermissionsRes, err error) {
	// 协作者列表属管理面信息，仅对管理者/所有者可见（RPC 无 operator，API 层校验）。
	if err := requireRole(l.ctx, l.svcCtx, req.ResourceType, req.ResourceId, req.UserId, "manager"); err != nil {
		return nil, err
	}
	res, err := l.svcCtx.Doc.ListPermissions(l.ctx, &doc_rpc.ListPermissionsReq{
		ResourceType: req.ResourceType,
		ResourceId:   req.ResourceId,
	})
	if err != nil {
		return nil, err
	}
	return &types.ListPermissionsRes{Permissions: convPermItems(res.Permissions)}, nil
}
