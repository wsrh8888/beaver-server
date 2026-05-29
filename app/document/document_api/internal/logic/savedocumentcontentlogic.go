package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/document/document_api/internal/svc"
	"beaver/app/document/document_api/internal/types"
	"beaver/app/document/document_models"
	"beaver/common/models/ctype"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SaveDocumentContentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveDocumentContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDocumentContentLogic {
	return &SaveDocumentContentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveDocumentContentLogic) SaveDocumentContent(req *types.SaveDocumentContentReq) (*types.SaveDocumentContentRes, error) {
	if _, err := validateDocumentFile(l.ctx, l.svcCtx.FileRpc, req.FileKey); err != nil {
		return nil, err
	}

	var doc document_models.CloudDocument
	err := l.svcCtx.DB.Where("doc_id = ? AND status = 1 AND deleted_at IS NULL", req.DocID).First(&doc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("文档不存在")
	}
	if err != nil {
		return nil, err
	}
	if doc.DocType == ctype.CloudDocTypeFolder {
		return nil, fmt.Errorf("文件夹不支持保存正文")
	}

	perm, err := resolveDocumentPerm(l.svcCtx.DB, &doc, req.UserID)
	if err != nil {
		return nil, err
	}
	if !canEdit(perm) {
		return nil, fmt.Errorf("无权限编辑该文档")
	}

	if req.Revision > 0 && req.Revision != doc.Revision {
		return nil, fmt.Errorf("文档已被他人更新，请刷新后重试")
	}

	var current document_models.CloudDocumentContent
	if err := l.svcCtx.DB.Where("doc_id = ?", doc.DocID).First(&current).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	newRevision := doc.Revision + 1
	if err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if current.DocID != "" {
			if err := tx.Create(&document_models.CloudDocumentRevision{
				DocID:    doc.DocID,
				Revision: doc.Revision,
				EditorID: req.UserID,
				FileKey:  current.FileKey,
			}).Error; err != nil {
				return err
			}
		}

		content := document_models.CloudDocumentContent{
			DocID:    doc.DocID,
			Revision: newRevision,
			FileKey:  req.FileKey,
		}
		if current.DocID == "" {
			if err := tx.Create(&content).Error; err != nil {
				return err
			}
		} else if err := tx.Model(&current).Updates(map[string]interface{}{
			"revision": newRevision,
			"file_key": req.FileKey,
		}).Error; err != nil {
			return err
		}

		return tx.Model(&doc).Updates(map[string]interface{}{
			"revision":       newRevision,
			"last_editor_id": req.UserID,
		}).Error
	}); err != nil {
		return nil, fmt.Errorf("保存失败: %w", err)
	}

	if err := l.svcCtx.DB.Where("doc_id = ? AND status = 1", req.DocID).First(&doc).Error; err != nil {
		return nil, err
	}

	return &types.SaveDocumentContentRes{
		Revision:  newRevision,
		UpdatedAt: doc.UpdatedAt.String(),
	}, nil
}
