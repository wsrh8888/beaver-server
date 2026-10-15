package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAgentAppsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAgentAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAgentAppsLogic {
	return &ListAgentAppsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAgentAppsLogic) ListAgentApps(in *platform_rpc.ListAgentAppsReq) (*platform_rpc.ListAgentAppsRes, error) {
	q := l.svcCtx.DB.Model(&agent_models.AgentApp{})
	if in.Scene != "" {
		q = q.Where("scene = ?", in.Scene)
	}
	if in.Status != 0 {
		q = q.Where("status = ?", int(in.Status))
	}
	var total int64
	q.Count(&total)
	page := int(in.Page)
	pageSize := int(in.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	var rows []agent_models.AgentApp
	q.Order("sort ASC, id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows)
	out := make([]*platform_rpc.AgentAppItem, 0, len(rows))
	for i := range rows {
		out = append(out, toAgentAppItem(&rows[i]))
	}
	return &platform_rpc.ListAgentAppsRes{Total: total, Apps: out}, nil
}
