package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAgentAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentAppLogic {
	return &CreateAgentAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAgentAppLogic) CreateAgentApp(in *platform_rpc.CreateAgentAppReq) (*platform_rpc.CreateAgentAppRes, error) {
	appID := "app_" + randHex(16)
	row := agent_models.AgentApp{
		AppID:           appID,
		Name:            in.Name,
		Scene:           in.Scene,
		SystemPrompt:    in.SystemPrompt,
		LLMModelID:       in.LlmModelId,
		ModelParams:      in.ModelParams,
		Tools:           in.Tools,
		McpServerIDs:    in.McpServerIds,
		KnowledgeIDs:    in.KnowledgeIds,
		OpeningMessage:   in.OpeningMessage,
		OpeningQuestions: in.OpeningQuestions,
		Variables:        in.Variables,
		Status:          int(in.Status),
		Sort:            int(in.Sort),
		CreatedBy:       in.CreatedBy,
	}
	if row.Status == 0 {
		row.Status = 1
	}
	if row.CreatedBy == "" {
		row.CreatedBy = "system"
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		return &platform_rpc.CreateAgentAppRes{}, nil
	}
	return &platform_rpc.CreateAgentAppRes{AppId: appID}, nil
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
