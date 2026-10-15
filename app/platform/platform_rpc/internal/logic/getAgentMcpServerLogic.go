package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentMcpServerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentMcpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentMcpServerLogic {
	return &GetAgentMcpServerLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// —— MCP 服务目录（agent/ 通用）——
func (l *GetAgentMcpServerLogic) GetAgentMcpServer(in *platform_rpc.GetAgentMcpServerReq) (*platform_rpc.GetAgentMcpServerRes, error) {
	var row agent_models.AgentMcpServer
	if err := l.svcCtx.DB.Where("mcp_server_id = ?", in.McpServerId).First(&row).Error; err != nil {
		return &platform_rpc.GetAgentMcpServerRes{}, nil
	}
	return &platform_rpc.GetAgentMcpServerRes{Server: toAgentMcpServerItem(&row)}, nil
}

func toAgentMcpServerItem(row *agent_models.AgentMcpServer) *platform_rpc.AgentMcpServerItem {
	if row == nil {
		return nil
	}
	return &platform_rpc.AgentMcpServerItem{
		Id:           uint64(row.Id),
		McpServerId:  row.McpServerID,
		Name:         row.Name,
		ServerName:   row.ServerName,
		Transport:    row.Transport,
		Config:       row.Config,
		Tools:        row.Tools,
		Status:       int32(row.Status),
	}
}
