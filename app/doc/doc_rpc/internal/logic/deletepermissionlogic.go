package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type DeletePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeletePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePermissionLogic {
	return &DeletePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DeletePermissionLogic", ctx),
	}
}

func (l *DeletePermissionLogic) DeletePermission(in *doc_rpc.DeletePermissionReq) (*doc_rpc.DeletePermissionRes, error) {
	rt := trim(in.ResourceType)
	if rt != "space" && rt != "doc" {
		return nil, errInvalidArg
	}
	opRole, err := effectiveRoleForResource(l.svcCtx.DB, rt, in.ResourceId, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canManage(opRole) {
		return nil, errNoPermission
	}

	pt := trim(in.PrincipalType)
	if pt == "" {
		pt = "user"
	}
	// 不删除 owner 行（归属权由转移流程管理）。
	if err := l.svcCtx.DB.Where("resource_type = ? AND resource_id = ? AND principal_type = ? AND principal_id = ? AND role <> ?",
		rt, in.ResourceId, pt, trim(in.PrincipalId), roleOwner).
		Delete(&doc_models.DocPermission{}).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "删除权限失败", Data: err})
		return nil, err
	}
	return &doc_rpc.DeletePermissionRes{}, nil
}
