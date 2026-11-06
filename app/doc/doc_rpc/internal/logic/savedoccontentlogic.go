package logic

import (
	"context"

	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
)

type SaveDocContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSaveDocContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDocContentLogic {
	return &SaveDocContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("SaveDocContentLogic", ctx),
	}
}

func (l *SaveDocContentLogic) SaveDocContent(in *doc_rpc.SaveDocContentReq) (*doc_rpc.SaveDocContentRes, error) {
	doc, err := loadDoc(l.svcCtx.DB, in.DocId)
	if err != nil {
		return nil, errNotFound
	}
	space, err := loadSpace(l.svcCtx.DB, doc.SpaceID)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveDocRole(l.svcCtx.DB, doc.SpaceID, doc.DocID, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canEdit(role) {
		return nil, errNoPermission
	}

	version, err := applyDocContent(l.ctx, l.svcCtx, doc, space, in.ContentJson, in.ContentText, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	return &doc_rpc.SaveDocContentRes{Version: int32(version)}, nil
}
