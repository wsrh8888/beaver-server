package logic

import (
	"context"
	"errors"
	"strings"
	"time"

	"beaver/app/agent/agent_api/internal/svc"
	"beaver/app/agent/agent_api/internal/types"
	"beaver/app/agent/agent_models"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAgentModelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAgentModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentModelLogic {
	return &CreateAgentModelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAgentModelLogic) CreateAgentModel(req *types.CreateAgentModelReq) (*types.CreateAgentModelRes, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return nil, errors.New("user id required")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("name required")
	}
	if strings.TrimSpace(req.Endpoint) == "" {
		return nil, errors.New("endpoint required")
	}
	if strings.TrimSpace(req.ApiKey) == "" {
		return nil, errors.New("api key required")
	}

	row := agent_models.AgentUserModel{
		ModelID:      "amodel_" + strings.ReplaceAll(uuid.NewString(), "-", ""),
		UserID:       strings.TrimSpace(req.UserID),
		Name:         strings.TrimSpace(req.Name),
		Endpoint:     strings.TrimSpace(req.Endpoint),
		ApiKey:       strings.TrimSpace(req.ApiKey),
		Tier:         normalizeModelTier(req.Tier),
		Capabilities: toModelCapabilities(req.Capabilities),
		Status:       1,
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		l.Errorf("CreateAgentModel failed: %v", err)
		return nil, err
	}

	return &types.CreateAgentModelRes{Model: toAgentModelItem(row)}, nil
}

func maskAgentModelApiKey(key string) string {
	runes := []rune(key)
	n := len(runes)
	if n == 0 {
		return ""
	}
	if n <= 8 {
		return "***"
	}
	return string(runes[:3]) + "***" + string(runes[n-4:])
}

func normalizeModelTier(tier string) string {
	switch strings.TrimSpace(tier) {
	case agent_models.ModelTierFast, agent_models.ModelTierStrong:
		return strings.TrimSpace(tier)
	default:
		return agent_models.ModelTierBalanced
	}
}

func toModelCapabilities(c types.AgentModelCapabilities) agent_models.AgentModelCapabilities {
	return agent_models.AgentModelCapabilities{
		Vision:           c.Vision,
		Tools:            c.Tools,
		Reasoning:        c.Reasoning,
		StructuredOutput: c.StructuredOutput,
		LongContext:      c.LongContext,
	}
}

func toAgentModelItem(row agent_models.AgentUserModel) types.AgentModelItem {
	return types.AgentModelItem{
		ModelId:      row.ModelID,
		Name:         row.Name,
		Endpoint:     row.Endpoint,
		ApiKeyMasked: maskAgentModelApiKey(row.ApiKey),
		Tier:         row.Tier,
		Capabilities: types.AgentModelCapabilities{
			Vision:           row.Capabilities.Vision,
			Tools:            row.Capabilities.Tools,
			Reasoning:        row.Capabilities.Reasoning,
			StructuredOutput: row.Capabilities.StructuredOutput,
			LongContext:      row.Capabilities.LongContext,
		},
		Status:    int32(row.Status),
		CreatedAt: time.Time(row.CreatedAt).UnixMilli(),
	}
}
