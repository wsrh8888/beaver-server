package logic

import (
	"context"

	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ResolveRoomEpochLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewResolveRoomEpochLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResolveRoomEpochLogic {
	return &ResolveRoomEpochLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ResolveRoomEpochLogic", ctx),
	}
}

func (l *ResolveRoomEpochLogic) ResolveRoomEpoch(in *doc_rpc.ResolveRoomEpochReq) (*doc_rpc.ResolveRoomEpochRes, error) {
	doc, err := loadDoc(l.svcCtx.DB, in.DocId)
	if err != nil {
		return &doc_rpc.ResolveRoomEpochRes{Valid: false}, nil
	}
	// epoch<=0 表示旧客户端尚未携带纪元；视为有效，但网关应使用返回值中的当前纪元建房间。
	valid := in.Epoch <= 0 || int(in.Epoch) == doc.Epoch
	return &doc_rpc.ResolveRoomEpochRes{Valid: valid, Epoch: int32(doc.Epoch)}, nil
}
