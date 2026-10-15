package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAgentToolLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentToolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentToolLogic {
	return &CreateAgentToolLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateAgentToolLogic) CreateAgentTool(in *platform_rpc.CreateAgentToolReq) (*platform_rpc.CreateAgentToolRes, error) {
	toolID := "tool_" + randHex(16)
	row := agent_models.AgentTool{
		ToolID:   toolID,
		Name:     in.Name,
		ToolType: in.ToolType,
		Config:   in.Config,
		Status:   int(in.Status),
		Sort:     int(in.Sort),
	}
	if row.ToolType == "" {
		row.ToolType = "builtin"
	}
	if row.Status == 0 {
		row.Status = 1
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		return &platform_rpc.CreateAgentToolRes{}, nil
	}
	return &platform_rpc.CreateAgentToolRes{ToolId: toolID}, nil
}
