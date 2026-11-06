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
	// 内部可信调用（doc_sync_ws 网关），权限在网关建连时按 token 校验。
	doc, err := loadDoc(l.svcCtx.DB, in.DocId)
	if err != nil {
		return nil, errNotFound
	}

	res := &doc_rpc.GetDocSyncStateRes{Version: int32(doc.Version)}

	var snap doc_models.DocYjsSnapshot
	upTo := int64(0)
	if err := l.svcCtx.DB.Where("doc_id = ?", doc.DocID).First(&snap).Error; err == nil {
		res.Snapshot = snap.State
		upTo = snap.UpdateSeqUpTo
	}

	// 快照之后的增量 update，按时序返回。
	var updates []doc_models.DocYjsUpdate
	l.svcCtx.DB.Where("doc_id = ? AND seq > ?", doc.DocID, upTo).Order("seq asc").Find(&updates)
	for i := range updates {
		res.Updates = append(res.Updates, updates[i].Update)
	}

	return res, nil
}
