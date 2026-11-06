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

type DeleteDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDocLogic {
	return &DeleteDocLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DeleteDocLogic", ctx),
	}
}

func (l *DeleteDocLogic) DeleteDoc(in *doc_rpc.DeleteDocReq) (*doc_rpc.DeleteDocRes, error) {
	doc, err := loadDoc(l.svcCtx.DB, in.DocId)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveDocRole(l.svcCtx.DB, doc.SpaceID, doc.DocID, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canManage(role) {
		return nil, errNoPermission
	}

	if err := l.svcCtx.DB.Model(&doc_models.DocDocument{}).Where("doc_id = ?", doc.DocID).Update("deleted_at", time.Now()).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "删除文档失败", Data: err})
		return nil, err
	}
	// 关联节点一并软删。
	if doc.NodeID != "" {
		_ = l.svcCtx.DB.Model(&doc_models.DocNode{}).Where("node_id = ?", doc.NodeID).Update("deleted_at", time.Now()).Error
	}
	deleteDocProjection(l.ctx, l.svcCtx, doc)

	return &doc_rpc.DeleteDocRes{}, nil
}
