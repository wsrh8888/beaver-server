package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ListNodesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListNodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNodesLogic {
	return &ListNodesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListNodesLogic", ctx),
	}
}

func (l *ListNodesLogic) ListNodes(in *doc_rpc.ListNodesReq) (*doc_rpc.ListNodesRes, error) {
	space, err := loadSpace(l.svcCtx.DB, in.SpaceId)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveSpaceRole(l.svcCtx.DB, space, in.UserId)
	if err != nil {
		return nil, err
	}
	if role == "" {
		return nil, errNoPermission
	}

	var rows []doc_models.DocNode
	if err := l.svcCtx.DB.Where("space_id = ? AND deleted_at = ?", space.SpaceID, "1970-01-01 00:00:00").
		Order("sort asc, id asc").Find(&rows).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "列出节点失败", Data: err})
		return nil, err
	}

	items := make([]*doc_rpc.NodeItem, 0, len(rows))
	for i := range rows {
		items = append(items, toNodeItem(&rows[i]))
	}
	return &doc_rpc.ListNodesRes{Nodes: items}, nil
}
