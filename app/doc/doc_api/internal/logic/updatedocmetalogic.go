package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type UpdateDocMetaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDocMetaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDocMetaLogic {
	return &UpdateDocMetaLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateDocMetaLogic) UpdateDocMeta(req *types.UpdateDocMetaReq) (resp *types.UpdateDocMetaRes, err error) {
	_, err = l.svcCtx.Doc.UpdateDocMeta(l.ctx, &doc_rpc.UpdateDocMetaReq{
		DocId:          req.DocId,
		Title:          req.Title,
		RagEnabled:     req.RagEnabled,
		Status:         int32(req.Status),
		OperatorUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.UpdateDocMetaRes{}, nil
}
