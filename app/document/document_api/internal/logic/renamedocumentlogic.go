package logic

import (
	"context"
	"fmt"

	"beaver/app/document/document_api/internal/svc"
	"beaver/app/document/document_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RenameDocumentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRenameDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RenameDocumentLogic {
	return &RenameDocumentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RenameDocumentLogic) RenameDocument(req *types.RenameDocumentReq) (*types.RenameDocumentRes, error) {
	doc, err := loadActiveDocument(l.svcCtx.DB, req.DocID)
	if err != nil {
		return nil, err
	}

	perm, err := resolveDocumentPerm(l.svcCtx.DB, doc, req.UserID)
	if err != nil {
		return nil, err
	}
	if !canEdit(perm) {
		return nil, fmt.Errorf("无权限重命名该文档")
	}

	if err := l.svcCtx.DB.Model(doc).Updates(map[string]interface{}{
		"title":          req.Title,
		"last_editor_id": req.UserID,
	}).Error; err != nil {
		return nil, fmt.Errorf("重命名失败: %w", err)
	}

	return &types.RenameDocumentRes{
		DocID: doc.DocID,
		Title: req.Title,
	}, nil
}
