package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type CreateTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTemplateLogic {
	return &CreateTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("CreateTemplateLogic", ctx),
	}
}

func (l *CreateTemplateLogic) CreateTemplate(in *doc_rpc.CreateTemplateReq) (*doc_rpc.CreateTemplateRes, error) {
	name := trim(in.Name)
	if name == "" {
		return nil, errInvalidArg
	}
	content := in.ContentJson
	if trim(content) == "" {
		content = skeletonContentJSON
	}
	// builtin 只由系统预置，用户创建一律降级为个人模板。
	scope := trim(in.Scope)
	if scope == "" || scope == "builtin" {
		scope = "owner"
	}
	dt := trim(in.DocType)
	if dt == "" {
		dt = "doc"
	}
	row := doc_models.DocTemplate{
		TemplateID:  genID("tpl"),
		Name:        name,
		Category:    trim(in.Category),
		Description: trim(in.Description),
		Cover:       trim(in.Cover),
		ContentJSON: content,
		DocType:     dt,
		Scope:       scope,
		OwnerUserID: trim(in.OperatorUserId),
		Status:      1,
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		return nil, err
	}
	auditJSON(l.svcCtx.DB, in.OperatorUserId, "template.create", "template", row.TemplateID,
		map[string]any{"name": name, "scope": scope})
	return &doc_rpc.CreateTemplateRes{TemplateId: row.TemplateID}, nil
}
