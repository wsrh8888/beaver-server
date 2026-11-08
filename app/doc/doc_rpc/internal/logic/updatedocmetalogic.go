package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type UpdateDocMetaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateDocMetaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDocMetaLogic {
	return &UpdateDocMetaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("UpdateDocMetaLogic", ctx),
	}
}

func (l *UpdateDocMetaLogic) UpdateDocMeta(in *doc_rpc.UpdateDocMetaReq) (*doc_rpc.UpdateDocMetaRes, error) {
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

	upd := map[string]any{}
	if title := trim(in.Title); title != "" {
		upd["title"] = title
	}
	if in.Status > 0 {
		upd["status"] = int(in.Status)
	}
	if err := l.svcCtx.DB.Model(&doc_models.DocDocument{}).Where("doc_id = ?", doc.DocID).Updates(upd).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "更新文档元信息失败", Data: err})
		return nil, err
	}

	if title := trim(in.Title); title != "" {
		doc.Title = title
	}
	if in.Status > 0 {
		doc.Status = int(in.Status)
	}

	return &doc_rpc.UpdateDocMetaRes{}, nil
}
