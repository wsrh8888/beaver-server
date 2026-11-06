package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type UpdateNodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNodeLogic {
	return &UpdateNodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("UpdateNodeLogic", ctx),
	}
}

func (l *UpdateNodeLogic) UpdateNode(in *doc_rpc.UpdateNodeReq) (*doc_rpc.UpdateNodeRes, error) {
	node, err := loadNode(l.svcCtx.DB, in.NodeId)
	if err != nil {
		return nil, errNotFound
	}
	space, err := loadSpace(l.svcCtx.DB, node.SpaceID)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveSpaceRole(l.svcCtx.DB, space, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canEdit(role) {
		return nil, errNoPermission
	}

	upd := map[string]any{}
	if title := trim(in.Title); title != "" {
		upd["title"] = title
	}
	if icon := trim(in.Icon); icon != "" {
		upd["icon"] = icon
	}
	if in.ParentId != node.ParentID {
		upd["parent_id"] = trim(in.ParentId)
	}
	if int(in.Sort) != node.Sort {
		upd["sort"] = int(in.Sort)
	}
	if len(upd) > 0 {
		if err := l.svcCtx.DB.Model(&doc_models.DocNode{}).Where("node_id = ?", node.NodeID).Updates(upd).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "更新节点失败", Data: err})
			return nil, err
		}
		// page 节点改名同步更新文档标题（元信息）。
		if title, ok := upd["title"].(string); ok && title != "" && node.DocID != "" {
			_ = l.svcCtx.DB.Model(&doc_models.DocDocument{}).Where("doc_id = ?", node.DocID).Update("title", title).Error
		}
	}

	return &doc_rpc.UpdateNodeRes{}, nil
}
