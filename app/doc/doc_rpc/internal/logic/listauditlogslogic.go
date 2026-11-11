package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ListAuditLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListAuditLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAuditLogsLogic {
	return &ListAuditLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListAuditLogsLogic", ctx),
	}
}

func (l *ListAuditLogsLogic) ListAuditLogs(in *doc_rpc.ListAuditLogsReq) (*doc_rpc.ListAuditLogsRes, error) {
	q := l.svcCtx.DB.Model(&doc_models.DocAuditLog{})
	if v := trim(in.ActorUserId); v != "" {
		q = q.Where("actor_user_id = ?", v)
	}
	if v := trim(in.Action); v != "" {
		q = q.Where("action = ?", v)
	}
	if v := trim(in.ResourceType); v != "" {
		q = q.Where("resource_type = ?", v)
	}
	if v := trim(in.ResourceId); v != "" {
		q = q.Where("resource_id = ?", v)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}
	page, size := pageWindow(in.Page, in.PageSize)
	var rows []doc_models.DocAuditLog
	if err := q.Order("created_at_str desc, id desc").
		Offset((page - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, err
	}
	logs := make([]*doc_rpc.AuditLogItem, 0, len(rows))
	for i := range rows {
		logs = append(logs, toAuditLogItem(&rows[i]))
	}
	return &doc_rpc.ListAuditLogsRes{Total: total, Logs: logs}, nil
}
