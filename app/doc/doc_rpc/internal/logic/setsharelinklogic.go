package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type SetShareLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSetShareLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetShareLinkLogic {
	return &SetShareLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("SetShareLinkLogic", ctx),
	}
}

func (l *SetShareLinkLogic) SetShareLink(in *doc_rpc.SetShareLinkReq) (*doc_rpc.SetShareLinkRes, error) {
	rt := trim(in.ResourceType)
	if rt != "doc" && rt != "space" {
		return nil, errInvalidArg
	}
	if rt == "doc" {
		if _, _, err := requireDoc(l.svcCtx.DB, in.ResourceId, in.OperatorUserId, roleManager); err != nil {
			return nil, err
		}
	} else {
		if _, _, err := requireSpace(l.svcCtx.DB, in.ResourceId, in.OperatorUserId, roleManager); err != nil {
			return nil, err
		}
	}
	// 链接不得授予管理权，避免一次分享造成越权提权。
	role := trim(in.Role)
	if roleRank(role) < roleRank(roleViewer) || roleRank(role) > roleRank(roleEditor) {
		return nil, errInvalidArg
	}
	scope := trim(in.Scope)
	if scope == "" {
		scope = "anyone"
	}
	if scope != "anyone" && scope != "org" && scope != "invited" {
		return nil, errInvalidArg
	}
	passwordHash := ""
	if p := trim(in.Password); p != "" {
		passwordHash = sha256hex(p)
	}

	var row doc_models.DocShareLink
	if err := l.svcCtx.DB.Where("resource_type = ? AND resource_id = ?", rt, in.ResourceId).
		First(&row).Error; err != nil {
		row = doc_models.DocShareLink{
			LinkID:       genID("lnk"),
			Token:        randHex(16),
			ResourceType: rt,
			ResourceID:   in.ResourceId,
			CreatedBy:    trim(in.OperatorUserId),
		}
	}
	if row.LinkID == "" {
		row.LinkID = genID("lnk")
	}
	if row.Token == "" {
		row.Token = randHex(16)
	}
	row.Scope = scope
	row.Role = role
	row.PasswordHash = passwordHash
	row.ExpiresAt = trim(in.ExpiresAt)
	row.AllowDownload = in.AllowDownload
	row.AllowCopy = in.AllowCopy
	row.AllowComment = in.AllowComment
	row.Watermark = in.Watermark
	row.Status = 1
	if err := l.svcCtx.DB.Save(&row).Error; err != nil {
		return nil, err
	}
	auditJSON(l.svcCtx.DB, in.OperatorUserId, "share.set_link", rt, in.ResourceId,
		map[string]any{"scope": scope, "role": role, "linkId": row.LinkID})
	return &doc_rpc.SetShareLinkRes{Link: toShareLinkItem(&row)}, nil
}
