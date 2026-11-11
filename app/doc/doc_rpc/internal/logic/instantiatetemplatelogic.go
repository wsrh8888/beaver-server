package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"

	"gorm.io/gorm"
)

type InstantiateTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewInstantiateTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstantiateTemplateLogic {
	return &InstantiateTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("InstantiateTemplateLogic", ctx),
	}
}

func (l *InstantiateTemplateLogic) InstantiateTemplate(in *doc_rpc.InstantiateTemplateReq) (*doc_rpc.InstantiateTemplateRes, error) {
	var tpl doc_models.DocTemplate
	if err := l.svcCtx.DB.Where("template_id = ? AND status = ?", in.TemplateId, 1).
		First(&tpl).Error; err != nil {
		return nil, errNotFound
	}
	actor := trim(in.OperatorUserId)
	if tpl.Scope == "owner" && trim(tpl.OwnerUserID) != actor {
		return nil, errNoPermission
	}
	space, _, err := requireSpace(l.svcCtx.DB, in.SpaceId, actor, roleEditor)
	if err != nil {
		return nil, err
	}
	title := trim(in.Title)
	if title == "" {
		title = tpl.Name
	}
	contentJSON := tpl.ContentJSON
	if trim(contentJSON) == "" {
		contentJSON = skeletonContentJSON
	}
	contentText := extractTextFromJSON(contentJSON)
	nodeID, docID := genID("node"), genID("doc")
	doc := doc_models.DocDocument{
		DocID:         docID,
		NodeID:        nodeID,
		SpaceID:       space.SpaceID,
		Title:         title,
		ContentJSON:   contentJSON,
		ContentText:   contentText,
		ContentDigest: sha256hex(contentText),
		Version:       1,
		Epoch:         1,
		DocType:       tpl.DocType,
		Cover:         tpl.Cover,
		OwnerUserID:   actor,
		Status:        1,
		DeletedAt:     sentinelDeletedAt,
	}
	node := doc_models.DocNode{
		NodeID:   nodeID,
		SpaceID:  space.SpaceID,
		ParentID: trim(in.ParentNodeId),
		Title:    title,
		Type:     "page",
		DocID:    docID,
	}
	if err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&node).Error; err != nil {
			return err
		}
		if err := tx.Create(&doc).Error; err != nil {
			return err
		}
		if err := tx.Model(&doc_models.DocTemplate{}).Where("template_id = ?", tpl.TemplateID).
			UpdateColumn("use_count", gorm.Expr("use_count + 1")).Error; err != nil {
			return err
		}
		auditJSON(tx, actor, "template.instantiate", "doc", docID,
			map[string]any{"templateId": tpl.TemplateID, "spaceId": space.SpaceID})
		return nil
	}); err != nil {
		return nil, err
	}
	return &doc_rpc.InstantiateTemplateRes{DocId: docID, NodeId: nodeID}, nil
}
