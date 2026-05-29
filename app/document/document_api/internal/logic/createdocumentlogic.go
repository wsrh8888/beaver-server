package logic

import (
	"context"
	"fmt"

	"beaver/app/document/document_api/internal/svc"
	"beaver/app/document/document_api/internal/types"
	"beaver/app/document/document_models"
	"beaver/common/models/ctype"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateDocumentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDocumentLogic {
	return &CreateDocumentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func defaultDocumentTitle(docType int) string {
	switch docType {
	case ctype.CloudDocTypeFolder:
		return "新建文件夹"
	case ctype.CloudDocTypeSheet:
		return "未命名表格"
	case ctype.CloudDocTypeSlide:
		return "未命名幻灯片"
	case ctype.CloudDocTypeMind:
		return "未命名思维笔记"
	default:
		return "未命名文档"
	}
}

func (l *CreateDocumentLogic) CreateDocument(req *types.CreateDocumentReq) (*types.CreateDocumentRes, error) {
	docType := ctype.CloudDocTypeDoc
	if req.DocType != nil {
		docType = *req.DocType
	}
	if docType < ctype.CloudDocTypeFolder || docType > ctype.CloudDocTypeMind {
		return nil, fmt.Errorf("不支持的文档类型")
	}
	if docType == ctype.CloudDocTypeFolder && req.FileKey != "" {
		return nil, fmt.Errorf("文件夹不能包含正文")
	}

	title := req.Title
	if title == "" {
		title = defaultDocumentTitle(docType)
	}

	if req.FileKey != "" {
		if _, err := validateDocumentFile(l.ctx, l.svcCtx.FileRpc, req.FileKey); err != nil {
			return nil, err
		}
	}

	docID := uuid.New().String()
	if err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		spaceID, err := resolveSpaceAccess(tx, req.UserID, req.SpaceID)
		if err != nil {
			return err
		}
		if err := checkParentWriteAccess(tx, req.UserID, spaceID, req.ParentID); err != nil {
			return err
		}

		doc := document_models.CloudDocument{
			DocID:        docID,
			OwnerID:      req.UserID,
			CreatorID:    req.UserID,
			DocType:      docType,
			Title:        title,
			Revision:     1,
			SpaceID:      spaceID,
			ParentID:     req.ParentID,
			LastEditorID: req.UserID,
			Status:       1,
		}
		if err := tx.Create(&doc).Error; err != nil {
			return err
		}
		if req.FileKey != "" {
			content := document_models.CloudDocumentContent{
				DocID:    docID,
				Revision: 1,
				FileKey:  req.FileKey,
			}
			if err := tx.Create(&content).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("创建失败: %w", err)
	}

	return &types.CreateDocumentRes{
		DocID:    docID,
		Title:    title,
		DocType:  docType,
		Revision: 1,
	}, nil
}
