package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type EmptyTrashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewEmptyTrashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmptyTrashLogic {
	return &EmptyTrashLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("EmptyTrashLogic", ctx),
	}
}

func (l *EmptyTrashLogic) EmptyTrash(in *doc_rpc.EmptyTrashReq) (*doc_rpc.EmptyTrashRes, error) {
	actor := trim(in.OperatorUserId)
	if actor == "" {
		return nil, errInvalidArg
	}
	var rows []doc_models.DocTrash
	if err := l.svcCtx.DB.Where("deleted_by = ?", actor).Find(&rows).Error; err != nil {
		return nil, err
	}
	purged := 0
	for i := range rows {
		if err := purgeOneTrash(l.svcCtx.DB, &rows[i]); err == nil {
			purged++
		}
	}
	auditJSON(l.svcCtx.DB, actor, "trash.empty", "trash", actor, map[string]any{"purged": purged})
	return &doc_rpc.EmptyTrashRes{Purged: int32(purged)}, nil
}
