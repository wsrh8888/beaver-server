package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type PurgeTrashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewPurgeTrashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PurgeTrashLogic {
	return &PurgeTrashLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("PurgeTrashLogic", ctx),
	}
}

func (l *PurgeTrashLogic) PurgeTrash(in *doc_rpc.PurgeTrashReq) (*doc_rpc.PurgeTrashRes, error) {
	actor := trim(in.OperatorUserId)
	var row doc_models.DocTrash
	if err := l.svcCtx.DB.Where("trash_id = ?", in.TrashId).First(&row).Error; err != nil {
		return nil, errNotFound
	}
	if row.DeletedBy != "" && row.DeletedBy != actor {
		return nil, errNoPermission
	}
	if err := purgeOneTrash(l.svcCtx.DB, &row); err != nil {
		return nil, err
	}
	auditJSON(l.svcCtx.DB, actor, "trash.purge", row.ResourceType, row.ResourceID,
		map[string]any{"trashId": row.TrashID})
	return &doc_rpc.PurgeTrashRes{}, nil
}
