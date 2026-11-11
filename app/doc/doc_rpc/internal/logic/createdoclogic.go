package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type CreateDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDocLogic {
	return &CreateDocLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("CreateDocLogic", ctx),
	}
}

func (l *CreateDocLogic) CreateDoc(in *doc_rpc.CreateDocReq) (*doc_rpc.CreateDocRes, error) {
	space, err := loadSpace(l.svcCtx.DB, in.SpaceId)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveSpaceRole(l.svcCtx.DB, space, in.OwnerUserId)
	if err != nil {
		return nil, err
	}
	if !canEdit(role) {
		return nil, errNoPermission
	}

	title := trim(in.Title)
	if title == "" {
		title = "未命名文档"
	}
	contentJSON := in.ContentJson
	if trim(contentJSON) == "" {
		contentJSON = skeletonContentJSON
	}

	docType := trim(in.DocType)
	if docType == "" {
		docType = "doc"
	}
	contentText := extractTextFromJSON(contentJSON)

	docID := genID("doc")
	row := doc_models.DocDocument{
		DocID:         docID,
		NodeID:        trim(in.NodeId),
		SpaceID:       space.SpaceID,
		Title:         title,
		ContentJSON:   contentJSON,
		ContentText:   contentText,
		ContentDigest: sha256hex(contentText),
		Version:       1,
		Epoch:         1,
		DocType:       docType,
		Icon:          trim(in.Icon),
		Cover:         trim(in.Cover),
		RagEnabled:    in.RagEnabled,
		OwnerUserID:   trim(in.OwnerUserId),
		Status:        1,
	}

	if err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
		if row.RagEnabled {
			enqueueIndex(tx, docID, row.Version, row.ContentDigest)
		}
		auditJSON(tx, in.OwnerUserId, "doc.create", "doc", docID, map[string]any{
			"title": title, "docType": docType, "spaceId": space.SpaceID,
		})
		return nil
	}); err != nil {
		l.logger.Error(model.LogMsg{Text: "创建文档失败", Data: err})
		return nil, err
	}

	return &doc_rpc.CreateDocRes{DocId: docID}, nil
}
