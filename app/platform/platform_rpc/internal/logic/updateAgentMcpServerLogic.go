package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAgentMcpServerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAgentMcpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentMcpServerLogic {
	return &UpdateAgentMcpServerLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateAgentMcpServerLogic) UpdateAgentMcpServer(in *platform_rpc.UpdateAgentMcpServerReq) (*platform_rpc.UpdateAgentMcpServerRes, error) {
	updates := map[string]any{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.ServerName != "" {
		updates["server_name"] = in.ServerName
	}
	if in.Transport != "" {
		updates["transport"] = in.Transport
	}
	if in.Config != "" {
		updates["config"] = in.Config
	}
	if in.Tools != "" {
		updates["tools"] = in.Tools
	}
	if in.Status != 0 {
		updates["status"] = int(in.Status)
	}
	if len(updates) > 0 {
		l.svcCtx.DB.Model(&agent_models.AgentMcpServer{}).Where("mcp_server_id = ?", in.McpServerId).Updates(updates)
	}
	return &platform_rpc.UpdateAgentMcpServerRes{}, nil
}
