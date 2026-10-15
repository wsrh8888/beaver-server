package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAgentToolLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAgentToolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentToolLogic {
	return &UpdateAgentToolLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateAgentToolLogic) UpdateAgentTool(in *platform_rpc.UpdateAgentToolReq) (*platform_rpc.UpdateAgentToolRes, error) {
	updates := map[string]any{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.ToolType != "" {
		updates["tool_type"] = in.ToolType
	}
	if in.Config != "" {
		updates["config"] = in.Config
	}
	if in.Status != 0 {
		updates["status"] = int(in.Status)
	}
	if in.Sort != 0 {
		updates["sort"] = int(in.Sort)
	}
	if len(updates) > 0 {
		l.svcCtx.DB.Model(&agent_models.AgentTool{}).Where("tool_id = ?", in.ToolId).Updates(updates)
	}
	return &platform_rpc.UpdateAgentToolRes{}, nil
}
