package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAgentAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAgentAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentAppLogic {
	return &UpdateAgentAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAgentAppLogic) UpdateAgentApp(in *platform_rpc.UpdateAgentAppReq) (*platform_rpc.UpdateAgentAppRes, error) {
	updates := map[string]any{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Scene != "" {
		updates["scene"] = in.Scene
	}
	if in.SystemPrompt != "" {
		updates["system_prompt"] = in.SystemPrompt
	}
	if in.LlmModelId != "" {
		updates["llm_model_id"] = in.LlmModelId
	}
	if in.ModelParams != "" {
		updates["model_params"] = in.ModelParams
	}
	if in.Tools != "" {
		updates["tools"] = in.Tools
	}
	if in.McpServerIds != "" {
		updates["mcp_server_ids"] = in.McpServerIds
	}
	if in.KnowledgeIds != "" {
		updates["knowledge_ids"] = in.KnowledgeIds
	}
	if in.OpeningMessage != "" {
		updates["opening_message"] = in.OpeningMessage
	}
	if in.OpeningQuestions != "" {
		updates["opening_questions"] = in.OpeningQuestions
	}
	if in.Variables != "" {
		updates["variables"] = in.Variables
	}
	if in.Status != 0 {
		updates["status"] = int(in.Status)
	}
	if in.Sort != 0 {
		updates["sort"] = int(in.Sort)
	}
	if len(updates) > 0 {
		l.svcCtx.DB.Model(&agent_models.AgentApp{}).Where("app_id = ?", in.AppId).Updates(updates)
	}
	return &platform_rpc.UpdateAgentAppRes{}, nil
}
