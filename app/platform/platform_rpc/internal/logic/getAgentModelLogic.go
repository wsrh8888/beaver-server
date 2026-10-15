package logic

import (
	"context"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentModelLogic {
	return &GetAgentModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// —— 模型目录（agent/ 通用）——
// source=custom：按 (model_id, user_id) 取用户自定义模型。
// source=official（默认）：按 model_id 取官方模型；"auto" 取 sort 最小启用官方模型。
func (l *GetAgentModelLogic) GetAgentModel(in *platform_rpc.GetAgentModelReq) (*platform_rpc.GetAgentModelRes, error) {
	source := in.ModelSource
	if source == "" {
		source = "official"
	}
	var row agent_models.AgentModel
	if source == "custom" {
		if err := l.svcCtx.DB.Where("source = 'custom' AND model_id = ? AND user_id = ? AND status = 1",
			in.ModelId, in.UserId).First(&row).Error; err != nil {
			return &platform_rpc.GetAgentModelRes{}, nil
		}
	} else {
		if in.ModelId == "" || in.ModelId == "auto" {
			if err := l.svcCtx.DB.Where("source = 'official' AND status = 1").
				Order("sort ASC, id ASC").First(&row).Error; err != nil {
				return &platform_rpc.GetAgentModelRes{}, nil
			}
		} else {
			if err := l.svcCtx.DB.Where("source = 'official' AND model_id = ? AND status = 1",
				in.ModelId).First(&row).Error; err != nil {
				return &platform_rpc.GetAgentModelRes{}, nil
			}
		}
	}
	return &platform_rpc.GetAgentModelRes{Model: toAgentModelItem(&row)}, nil
}

func toAgentModelItem(row *agent_models.AgentModel) *platform_rpc.AgentModelItem {
	if row == nil {
		return nil
	}
	var tiers agent_models.TokenTiers
	if row.TokenTiers != nil {
		tiers = *row.TokenTiers
	}
	item := &platform_rpc.AgentModelItem{
		Id:                uint64(row.Id),
		ModelId:           row.ModelID,
		Source:            row.Source,
		Name:              row.Name,
		ModelName:         strPtr(row.ModelName),
		Endpoint:          strPtr(row.Endpoint),
		Provider:          strPtr(row.Provider),
		ContextWindow:     row.ContextWindow,
		MaxOutputTokens:   int32(row.MaxOutputTokens),
		MaxTokensFlash:    int32(tiers.Flash),
		MaxTokensLow:      int32(tiers.Low),
		MaxTokensMedium:   int32(tiers.Medium),
		MaxTokensHigh:     int32(tiers.High),
		InstructionRole:   row.InstructionRole,
		SupportsThinking:  row.SupportsThinking,
		ThinkingCanDisable: row.ThinkingCanDisable,
		SupportsTools:     row.SupportsTools,
		SupportsVision:   row.SupportsVision,
		Sort:              int32(row.Sort),
		Status:            int32(row.Status),
		InputModalities:   row.InputModalities,
		OutputModalities:  row.OutputModalities,
		ThinkingLevels:    row.ThinkingLevels,
	}
	if row.UserID != nil {
		item.UserId = *row.UserID
	}
	if row.APIKey != nil {
		// agent 跑 LLM 需要真实 api_key；后台管理走 platform_api HTTP 接口（那个掩码）
		item.ApiKey = *row.APIKey
	}
	return item
}

func strPtr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
