package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ResolveShareLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewResolveShareLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResolveShareLinkLogic {
	return &ResolveShareLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ResolveShareLinkLogic", ctx),
	}
}

func (l *ResolveShareLinkLogic) ResolveShareLink(in *doc_rpc.ResolveShareLinkReq) (*doc_rpc.ResolveShareLinkRes, error) {
	var row doc_models.DocShareLink
	if err := l.svcCtx.DB.Where("token = ? AND status = ?", trim(in.Token), 1).
		First(&row).Error; err != nil {
		return nil, errNotFound
	}
	if linkExpired(row.ExpiresAt) {
		return nil, errNotFound
	}
	if row.PasswordHash != "" && sha256hex(trim(in.Password)) != row.PasswordHash {
		return nil, errNoPermission
	}
	// invited 链接要求访问者具备登录身份；后续读写仍走同一套 CheckPermission。
	if row.Scope == "invited" && trim(in.UserId) == "" {
		return nil, errNoPermission
	}

	res := &doc_rpc.ResolveShareLinkRes{
		ResourceType:  row.ResourceType,
		ResourceId:    row.ResourceID,
		Role:          row.Role,
		AllowDownload: row.AllowDownload,
		AllowCopy:     row.AllowCopy,
		AllowComment:  row.AllowComment,
		Watermark:     row.Watermark,
	}
	if row.ResourceType == "doc" {
		doc, err := loadDoc(l.svcCtx.DB, row.ResourceID)
		if err != nil {
			return nil, errNotFound
		}
		res.DocId = doc.DocID
		res.Title = doc.Title
	} else {
		space, err := loadSpace(l.svcCtx.DB, row.ResourceID)
		if err != nil {
			return nil, errNotFound
		}
		res.Title = space.Name
	}
	_ = l.svcCtx.DB.Model(&doc_models.DocShareLink{}).
		Where("link_id = ?", row.LinkID).
		Update("view_count", row.ViewCount+1).Error
	return res, nil
}
