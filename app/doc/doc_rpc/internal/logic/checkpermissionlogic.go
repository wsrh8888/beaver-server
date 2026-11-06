package logic

import (
	"context"

	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
)

type CheckPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCheckPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPermissionLogic {
	return &CheckPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("CheckPermissionLogic", ctx),
	}
}

func (l *CheckPermissionLogic) CheckPermission(in *doc_rpc.CheckPermissionReq) (*doc_rpc.CheckPermissionRes, error) {
	rt := trim(in.ResourceType)
	if rt != "space" && rt != "doc" {
		return nil, errInvalidArg
	}
	// 无权限返回空串而非错误：这是查询语义，调用方据 role 判定。
	role, err := effectiveRoleForResource(l.svcCtx.DB, rt, in.ResourceId, in.UserId)
	if err != nil {
		if err == errNotFound {
			return &doc_rpc.CheckPermissionRes{Role: ""}, nil
		}
		return nil, err
	}
	return &doc_rpc.CheckPermissionRes{Role: role}, nil
}
