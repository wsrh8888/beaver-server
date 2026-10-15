package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAgentModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentModelLogic {
	return &CreateAgentModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAgentModelLogic) CreateAgentModel(in *platform_rpc.CreateAgentModelReq) (*platform_rpc.CreateAgentModelRes, error) {
	modelID := "amodel_" + randHex(16)
	source := in.Source
	if source == "" {
		source = "custom"
	}
	var userID *string
	if in.UserId != "" {
		uid := in.UserId
		userID = &uid
	}
	var provider *string
	if in.Provider != "" {
		p := in.Provider
		provider = &p
	}
	var modelName *string
	if in.ModelName != "" {
		mn := in.ModelName
		modelName = &mn
	}
	var endpoint *string
	if in.Endpoint != "" {
		ep := in.Endpoint
		endpoint = &ep
	}
	var apiKey *string
	if in.ApiKey != "" {
		ak := in.ApiKey
		apiKey = &ak
	}
	status := int(in.Status)
	if status == 0 {
		status = 1
	}
	row := agent_models.AgentModel{
		ModelID:          modelID,
		Source:           source,
		UserID:           userID,
		Name:             in.Name,
		Provider:         provider,
		ModelName:        modelName,
		Endpoint:         endpoint,
		APIKey:           apiKey,
		ContextWindow:    in.ContextWindow,
		MaxOutputTokens:  int(in.MaxOutputTokens),
		InputModalities:  in.InputModalities,
		OutputModalities: in.OutputModalities,
		TokenTiers: &agent_models.TokenTiers{
			Flash:  int(in.MaxTokensFlash),
			Low:    int(in.MaxTokensLow),
			Medium: int(in.MaxTokensMedium),
			High:   int(in.MaxTokensHigh),
		},
		InstructionRole:  in.InstructionRole,
		SupportsThinking: in.SupportsThinking,
		ThinkingCanDisable: in.ThinkingCanDisable,
		ThinkingLevels:   in.ThinkingLevels,
		SupportsTools:    in.SupportsTools,
		SupportsVision:   in.SupportsVision,
		Sort:             int(in.Sort),
		Status:           status,
	}
	if row.InstructionRole == "" {
		row.InstructionRole = "system"
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		return &platform_rpc.CreateAgentModelRes{}, nil
	}
	return &platform_rpc.CreateAgentModelRes{ModelId: modelID}, nil
}
