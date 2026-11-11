package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type GetShareLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetShareLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShareLinkLogic {
	return &GetShareLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetShareLinkLogic", ctx),
	}
}

func (l *GetShareLinkLogic) GetShareLink(in *doc_rpc.GetShareLinkReq) (*doc_rpc.GetShareLinkRes, error) {
	rt := trim(in.ResourceType)
	if rt == "doc" {
		if _, _, err := requireDoc(l.svcCtx.DB, in.ResourceId, in.OperatorUserId, roleViewer); err != nil {
			return nil, err
		}
	} else if rt == "space" {
		if _, _, err := requireSpace(l.svcCtx.DB, in.ResourceId, in.OperatorUserId, roleViewer); err != nil {
			return nil, err
		}
	} else {
		return nil, errInvalidArg
	}
	var row doc_models.DocShareLink
	if err := l.svcCtx.DB.Where("resource_type = ? AND resource_id = ? AND status = ?",
		rt, in.ResourceId, 1).First(&row).Error; err != nil {
		return &doc_rpc.GetShareLinkRes{}, nil
	}
	return &doc_rpc.GetShareLinkRes{Link: toShareLinkItem(&row)}, nil
}
