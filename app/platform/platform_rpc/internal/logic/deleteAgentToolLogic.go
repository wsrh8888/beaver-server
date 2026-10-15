package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAgentToolLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAgentToolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentToolLogic {
	return &DeleteAgentToolLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteAgentToolLogic) DeleteAgentTool(in *platform_rpc.DeleteAgentToolReq) (*platform_rpc.DeleteAgentToolRes, error) {
	l.svcCtx.DB.Where("tool_id = ?", in.ToolId).Delete(&agent_models.AgentTool{})
	return &platform_rpc.DeleteAgentToolRes{}, nil
}
