package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAgentMcpServerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAgentMcpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentMcpServerLogic {
	return &DeleteAgentMcpServerLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteAgentMcpServerLogic) DeleteAgentMcpServer(in *platform_rpc.DeleteAgentMcpServerReq) (*platform_rpc.DeleteAgentMcpServerRes, error) {
	l.svcCtx.DB.Where("mcp_server_id = ?", in.McpServerId).Delete(&agent_models.AgentMcpServer{})
	return &platform_rpc.DeleteAgentMcpServerRes{}, nil
}
