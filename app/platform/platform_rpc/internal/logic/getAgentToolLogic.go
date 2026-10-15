package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentToolLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentToolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentToolLogic {
	return &GetAgentToolLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetAgentToolLogic) GetAgentTool(in *platform_rpc.GetAgentToolReq) (*platform_rpc.GetAgentToolRes, error) {
	var row agent_models.AgentTool
	if err := l.svcCtx.DB.Where("tool_id = ?", in.ToolId).First(&row).Error; err != nil {
		return &platform_rpc.GetAgentToolRes{}, nil
	}
	return &platform_rpc.GetAgentToolRes{Tool: toAgentToolItem(&row)}, nil
}

func toAgentToolItem(row *agent_models.AgentTool) *platform_rpc.AgentToolItem {
	if row == nil {
		return nil
	}
	return &platform_rpc.AgentToolItem{
		Id:       uint64(row.Id),
		ToolId:   row.ToolID,
		Name:     row.Name,
		ToolType: row.ToolType,
		Config:   row.Config,
		Status:   int32(row.Status),
		Sort:     int32(row.Sort),
	}
}
