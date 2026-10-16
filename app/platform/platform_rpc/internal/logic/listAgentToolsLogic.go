package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAgentToolsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAgentToolsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAgentToolsLogic {
	return &ListAgentToolsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// —— 工具目录（agent/ 通用）——
func (l *ListAgentToolsLogic) ListAgentTools(in *platform_rpc.ListAgentToolsReq) (*platform_rpc.ListAgentToolsRes, error) {
	q := l.svcCtx.DB.Model(&agent_models.AgentTool{})
	if in.Status != 0 {
		q = q.Where("status = ?", int(in.Status))
	}
	if in.Source != "" {
		q = q.Where("source = ?", in.Source)
	}
	var rows []agent_models.AgentTool
	q.Order("sort ASC, id ASC").Find(&rows)
	out := make([]*platform_rpc.AgentToolItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, &platform_rpc.AgentToolItem{
			Id:       uint64(row.Id),
			ToolId:   row.ToolID,
			Name:     row.Name,
			ToolType: row.ToolType,
			Config:   row.Config,
			Status:   int32(row.Status),
			Sort:     int32(row.Sort),
			Source:   row.Source,
			UserId:   row.UserID,
		})
	}
	return &platform_rpc.ListAgentToolsRes{Tools: out}, nil
}
