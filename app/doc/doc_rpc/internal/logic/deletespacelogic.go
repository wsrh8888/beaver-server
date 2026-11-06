package logic

import (
	"context"
	"time"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type DeleteSpaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSpaceLogic {
	return &DeleteSpaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DeleteSpaceLogic", ctx),
	}
}

func (l *DeleteSpaceLogic) DeleteSpace(in *doc_rpc.DeleteSpaceReq) (*doc_rpc.DeleteSpaceRes, error) {
	space, err := loadSpace(l.svcCtx.DB, in.SpaceId)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveSpaceRole(l.svcCtx.DB, space, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canManage(role) {
		return nil, errNoPermission
	}

	now := time.Now()
	if err := l.svcCtx.DB.Model(&doc_models.DocSpace{}).Where("space_id = ?", space.SpaceID).Update("deleted_at", now).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "删除空间失败", Data: err})
		return nil, err
	}
	// 级联软删其下节点与文档。
	if err := l.svcCtx.DB.Model(&doc_models.DocNode{}).Where("space_id = ? AND deleted_at = ?", space.SpaceID, "1970-01-01 00:00:00").Update("deleted_at", now).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "级联删除节点失败", Data: err})
	}
	var docs []doc_models.DocDocument
	l.svcCtx.DB.Where("space_id = ? AND deleted_at = ?", space.SpaceID, "1970-01-01 00:00:00").Find(&docs)
	for i := range docs {
		if err := l.svcCtx.DB.Model(&doc_models.DocDocument{}).Where("doc_id = ?", docs[i].DocID).Update("deleted_at", now).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "级联删除文档失败", Data: err})
		}
	}
	// 清 RAG 投影：整库删除 collection（级联软删其下 knowledge document 并清索引）。
	deleteSpaceProjection(l.ctx, l.svcCtx, space)

	return &doc_rpc.DeleteSpaceRes{}, nil
}
