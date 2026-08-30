package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/agent/agent_api/internal/svc"
	"beaver/app/agent/agent_api/internal/types"
	"beaver/app/agent/agent_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateAgentModelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAgentModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentModelLogic {
	return &UpdateAgentModelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAgentModelLogic) UpdateAgentModel(req *types.UpdateAgentModelReq) (*types.UpdateAgentModelRes, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return nil, errors.New("user id required")
	}
	if strings.TrimSpace(req.ModelId) == "" {
		return nil, errors.New("modelId required")
	}

	var row agent_models.AgentUserModel
	err := l.svcCtx.DB.Where("model_id = ? AND user_id = ?", strings.TrimSpace(req.ModelId), strings.TrimSpace(req.UserID)).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("model not found")
	}
	if err != nil {
		l.Errorf("UpdateAgentModel query failed: %v", err)
		return nil, err
	}

	if name := strings.TrimSpace(req.Name); name != "" {
		row.Name = name
	}
	if endpoint := strings.TrimSpace(req.Endpoint); endpoint != "" {
		row.Endpoint = endpoint
	}
	if apiKey := strings.TrimSpace(req.ApiKey); apiKey != "" {
		row.ApiKey = apiKey
	}
	if tier := strings.TrimSpace(req.Tier); tier != "" {
		row.Tier = normalizeModelTier(tier)
	}
	row.Capabilities = toModelCapabilities(req.Capabilities)

	if err := l.svcCtx.DB.Save(&row).Error; err != nil {
		l.Errorf("UpdateAgentModel save failed: %v", err)
		return nil, err
	}

	return &types.UpdateAgentModelRes{Model: toAgentModelItem(row)}, nil
}
