package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
)

type SaveDocContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveDocContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDocContentLogic {
	return &SaveDocContentLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *SaveDocContentLogic) SaveDocContent(req *types.SaveDocContentReq) (resp *types.SaveDocContentRes, err error) {
	res, err := l.svcCtx.Doc.SaveDocContent(l.ctx, &doc_rpc.SaveDocContentReq{
		DocId:          req.DocId,
		ContentJson:    req.ContentJson,
		ContentText:    req.ContentText,
		OperatorUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.SaveDocContentRes{Version: int(res.Version)}, nil
}
