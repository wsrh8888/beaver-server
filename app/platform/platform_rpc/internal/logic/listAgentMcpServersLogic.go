package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAgentMcpServersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAgentMcpServersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAgentMcpServersLogic {
	return &ListAgentMcpServersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// —— MCP 服务目录（agent/ 通用）——
func (l *ListAgentMcpServersLogic) ListAgentMcpServers(in *platform_rpc.ListAgentMcpServersReq) (*platform_rpc.ListAgentMcpServersRes, error) {
	q := l.svcCtx.DB.Model(&agent_models.AgentMcpServer{})
	if in.Status != 0 {
		q = q.Where("status = ?", int(in.Status))
	}
	if in.Source != "" {
		q = q.Where("source = ?", in.Source)
	}
	var rows []agent_models.AgentMcpServer
	q.Order("id ASC").Find(&rows)
	out := make([]*platform_rpc.AgentMcpServerItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, &platform_rpc.AgentMcpServerItem{
			Id:         uint64(row.Id),
			McpServerId: row.McpServerID,
			Name:       row.Name,
			ServerName: row.ServerName,
			Transport:  row.Transport,
			Config:     row.Config,
			Tools:      row.Tools,
			Status:     int32(row.Status),
			Source:     row.Source,
			UserId:     row.UserID,
		})
	}
	return &platform_rpc.ListAgentMcpServersRes{Servers: out}, nil
}
