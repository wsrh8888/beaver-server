package logic

import (
	"context"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"

	"github.com/zeromicro/go-zero/core/logx"

	beaverlog "beaver/utils/beaverlog"
)

// goctl 会把 logx 写进 imports，占位避免未使用导入。
var _ = logx.WithContext

type ReindexDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 触发文档 RAG 重索引
func NewReindexDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReindexDocLogic {
	return &ReindexDocLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ReindexDoc", ctx),
	}
}

func (l *ReindexDocLogic) ReindexDoc(req *types.ReindexDocReq) (resp *types.ReindexDocRes, err error) {
	res, err := l.svcCtx.Doc.ReindexDoc(l.ctx, &doc_rpc.ReindexDocReq{
		DocId:          req.DocId,
		OperatorUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.ReindexDocRes{Ok: res.Ok, KnowledgeDocId: res.KnowledgeDocId}, nil
}
