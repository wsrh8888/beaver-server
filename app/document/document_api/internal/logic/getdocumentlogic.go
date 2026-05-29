package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/document/document_api/internal/svc"
	"beaver/app/document/document_api/internal/types"
	"beaver/app/document/document_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetDocumentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDocumentLogic {
	return &GetDocumentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDocumentLogic) GetDocument(req *types.GetDocumentReq) (*types.GetDocumentRes, error) {
	var doc document_models.CloudDocument
	err := l.svcCtx.DB.Where("doc_id = ? AND status = 1 AND deleted_at IS NULL", req.DocID).First(&doc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("文档不存在")
	}
	if err != nil {
		return nil, err
	}

	perm, err := resolveDocumentPerm(l.svcCtx.DB, &doc, req.UserID)
	if err != nil {
		return nil, err
	}
	if !canRead(perm) {
		return nil, fmt.Errorf("无权限访问该文档")
	}

	res := &types.GetDocumentRes{
		DocID:     doc.DocID,
		OwnerID:   doc.OwnerID,
		Title:     doc.Title,
		DocType:   doc.DocType,
		Revision:  doc.Revision,
		Perm:      perm,
		UpdatedAt: doc.UpdatedAt.String(),
	}

	var body document_models.CloudDocumentContent
	if err := l.svcCtx.DB.Where("doc_id = ?", doc.DocID).First(&body).Error; err == nil {
		res.FileKey = body.FileKey
		if detail, err := validateDocumentFile(l.ctx, l.svcCtx.FileRpc, body.FileKey); err == nil {
			res.FileSize = detail.Size
			res.FileMd5 = detail.Md5
		}
	}

	return res, nil
}
