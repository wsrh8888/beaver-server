package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type DeleteTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTemplateLogic {
	return &DeleteTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DeleteTemplateLogic", ctx),
	}
}

func (l *DeleteTemplateLogic) DeleteTemplate(in *doc_rpc.DeleteTemplateReq) (*doc_rpc.DeleteTemplateRes, error) {
	var tpl doc_models.DocTemplate
	if err := l.svcCtx.DB.Where("template_id = ?", in.TemplateId).First(&tpl).Error; err != nil {
		return nil, errNotFound
	}
	actor := trim(in.OperatorUserId)
	if tpl.Scope == "builtin" || (trim(tpl.OwnerUserID) != "" && trim(tpl.OwnerUserID) != actor) {
		return nil, errNoPermission
	}
	if err := l.svcCtx.DB.Model(&doc_models.DocTemplate{}).
		Where("template_id = ?", tpl.TemplateID).Update("status", 2).Error; err != nil {
		return nil, err
	}
	auditJSON(l.svcCtx.DB, actor, "template.delete", "template", tpl.TemplateID,
		map[string]any{"name": tpl.Name})
	return &doc_rpc.DeleteTemplateRes{}, nil
}
