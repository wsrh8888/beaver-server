package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAgentMcpServerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentMcpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentMcpServerLogic {
	return &CreateAgentMcpServerLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateAgentMcpServerLogic) CreateAgentMcpServer(in *platform_rpc.CreateAgentMcpServerReq) (*platform_rpc.CreateAgentMcpServerRes, error) {
	mcpServerID := "mcp_" + randHex(16)
	row := agent_models.AgentMcpServer{
		McpServerID: mcpServerID,
		Name:        in.Name,
		ServerName:  in.ServerName,
		Transport:   in.Transport,
		Config:      in.Config,
		Tools:       in.Tools,
		Status:      int(in.Status),
	}
	if row.Transport == "" {
		row.Transport = "stdio"
	}
	if row.Status == 0 {
		row.Status = 1
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		return &platform_rpc.CreateAgentMcpServerRes{}, nil
	}
	return &platform_rpc.CreateAgentMcpServerRes{McpServerId: mcpServerID}, nil
}
