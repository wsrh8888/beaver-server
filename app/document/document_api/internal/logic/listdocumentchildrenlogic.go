package logic

import (
	"context"

	"beaver/app/document/document_api/internal/svc"
	"beaver/app/document/document_api/internal/types"
	"beaver/app/document/document_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDocumentChildrenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDocumentChildrenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDocumentChildrenLogic {
	return &ListDocumentChildrenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDocumentChildrenLogic) ListDocumentChildren(req *types.ListDocumentChildrenReq) (*types.ListDocumentChildrenRes, error) {
	spaceID, err := resolveSpaceAccess(l.svcCtx.DB, req.UserID, req.SpaceID)
	if err != nil {
		return nil, err
	}
	if err := checkParentReadAccess(l.svcCtx.DB, req.UserID, spaceID, req.ParentID); err != nil {
		return nil, err
	}

	var docs []document_models.CloudDocument
	if err := l.svcCtx.DB.Where("space_id = ? AND parent_id = ? AND status = 1 AND deleted_at IS NULL",
		spaceID, req.ParentID).Order("sort_order asc, id asc").Find(&docs).Error; err != nil {
		return nil, err
	}

	list := make([]types.ListDocumentChildrenItem, 0, len(docs))
	for i := range docs {
		perm, err := resolveDocumentPerm(l.svcCtx.DB, &docs[i], req.UserID)
		if err != nil {
			return nil, err
		}
		if !canRead(perm) {
			continue
		}
		list = append(list, types.ListDocumentChildrenItem{
			DocID:     docs[i].DocID,
			Title:     docs[i].Title,
			DocType:   docs[i].DocType,
			SortOrder: docs[i].SortOrder,
			UpdatedAt: docs[i].UpdatedAt.String(),
		})
	}

	return &types.ListDocumentChildrenRes{
		SpaceID:  spaceID,
		ParentID: req.ParentID,
		List:     list,
	}, nil
}
