package logic

import (
	"context"
	"fmt"

	"beaver/app/document/document_api/internal/svc"
	"beaver/app/document/document_api/internal/types"
	"beaver/common/models/ctype"

	"github.com/zeromicro/go-zero/core/logx"
)

type MoveDocumentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMoveDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MoveDocumentLogic {
	return &MoveDocumentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MoveDocumentLogic) MoveDocument(req *types.MoveDocumentReq) (*types.MoveDocumentRes, error) {
	doc, err := loadActiveDocument(l.svcCtx.DB, req.DocID)
	if err != nil {
		return nil, err
	}

	perm, err := resolveDocumentPerm(l.svcCtx.DB, doc, req.UserID)
	if err != nil {
		return nil, err
	}
	if !canEdit(perm) {
		return nil, fmt.Errorf("无权限移动该文档")
	}

	destParentID := req.ParentID
	if destParentID == doc.DocID {
		return nil, fmt.Errorf("不能移动到自身目录下")
	}
	if destParentID == doc.ParentID {
		return &types.MoveDocumentRes{
			DocID:    doc.DocID,
			ParentID: doc.ParentID,
		}, nil
	}

	if err := checkParentWriteAccess(l.svcCtx.DB, req.UserID, doc.SpaceID, destParentID); err != nil {
		return nil, err
	}

	if doc.DocType == ctype.CloudDocTypeFolder && destParentID != "" {
		isDesc, err := isDescendantOf(l.svcCtx.DB, doc.DocID, destParentID)
		if err != nil {
			return nil, err
		}
		if isDesc {
			return nil, fmt.Errorf("不能移动到自身子目录下")
		}
	}

	if err := l.svcCtx.DB.Model(doc).Updates(map[string]interface{}{
		"parent_id":      destParentID,
		"last_editor_id": req.UserID,
	}).Error; err != nil {
		return nil, fmt.Errorf("移动失败: %w", err)
	}

	return &types.MoveDocumentRes{
		DocID:    doc.DocID,
		ParentID: destParentID,
	}, nil
}
