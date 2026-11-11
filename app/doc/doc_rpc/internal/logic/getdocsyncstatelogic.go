package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
)

type GetDocSyncStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetDocSyncStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDocSyncStateLogic {
	return &GetDocSyncStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetDocSyncStateLogic", ctx),
	}
}

func (l *GetDocSyncStateLogic) GetDocSyncState(in *doc_rpc.GetDocSyncStateReq) (*doc_rpc.GetDocSyncStateRes, error) {
	// 内部可信调用（doc_api 协同网关），权限在网关建连时按 token 校验。
	doc, err := loadDoc(l.svcCtx.DB, in.DocId)
	if err != nil {
		return nil, errNotFound
	}

	// 始终返回文档当前纪元。请求携带旧纪元时不下发任何 CRDT 状态，
	// 迫使旧房间客户端以空状态重建到新纪元，避免回滚后旧状态被写回。
	res := &doc_rpc.GetDocSyncStateRes{
		Version: int32(doc.Version),
		Epoch:   int32(doc.Epoch),
	}
	if in.Epoch > 0 && int(in.Epoch) != doc.Epoch {
		return res, nil
	}

	var snap doc_models.DocYjsSnapshot
	upTo := int64(0)
	if err := l.svcCtx.DB.Where("doc_id = ? AND epoch = ?", doc.DocID, doc.Epoch).First(&snap).Error; err == nil {
		res.Snapshot = snap.State
		upTo = snap.UpdateSeqUpTo
	}

	// 快照之后的增量 update，严格限定在当前纪元内，按时序返回。
	var updates []doc_models.DocYjsUpdate
	l.svcCtx.DB.Where("doc_id = ? AND epoch = ? AND seq > ?", doc.DocID, doc.Epoch, upTo).
		Order("seq asc").Find(&updates)
	for i := range updates {
		res.Updates = append(res.Updates, updates[i].Update)
	}

	return res, nil
}
