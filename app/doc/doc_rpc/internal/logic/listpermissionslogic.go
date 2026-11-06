package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ListPermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPermissionsLogic {
	return &ListPermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListPermissionsLogic", ctx),
	}
}

func (l *ListPermissionsLogic) ListPermissions(in *doc_rpc.ListPermissionsReq) (*doc_rpc.ListPermissionsRes, error) {
	var rows []doc_models.DocPermission
	if err := l.svcCtx.DB.Where("resource_type = ? AND resource_id = ?", trim(in.ResourceType), in.ResourceId).
		Order("id asc").Find(&rows).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "列出权限失败", Data: err})
		return nil, err
	}
	items := make([]*doc_rpc.PermItem, 0, len(rows))
	for i := range rows {
		items = append(items, toPermItem(&rows[i]))
	}
	return &doc_rpc.ListPermissionsRes{Permissions: items}, nil
}
