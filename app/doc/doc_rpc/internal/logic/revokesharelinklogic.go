package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type RevokeShareLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRevokeShareLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeShareLinkLogic {
	return &RevokeShareLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("RevokeShareLinkLogic", ctx),
	}
}

func (l *RevokeShareLinkLogic) RevokeShareLink(in *doc_rpc.RevokeShareLinkReq) (*doc_rpc.RevokeShareLinkRes, error) {
	var row doc_models.DocShareLink
	if err := l.svcCtx.DB.Where("link_id = ?", in.LinkId).First(&row).Error; err != nil {
		return nil, errNotFound
	}
	if row.ResourceType == "doc" {
		if _, _, err := requireDoc(l.svcCtx.DB, row.ResourceID, in.OperatorUserId, roleManager); err != nil {
			return nil, err
		}
	} else {
		if _, _, err := requireSpace(l.svcCtx.DB, row.ResourceID, in.OperatorUserId, roleManager); err != nil {
			return nil, err
		}
	}
	if err := l.svcCtx.DB.Model(&doc_models.DocShareLink{}).
		Where("link_id = ?", row.LinkID).Update("status", 2).Error; err != nil {
		return nil, err
	}
	auditJSON(l.svcCtx.DB, in.OperatorUserId, "share.revoke_link", row.ResourceType, row.ResourceID,
		map[string]any{"linkId": row.LinkID})
	return &doc_rpc.RevokeShareLinkRes{}, nil
}
