package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/agent/agent_api/internal/svc"
	"beaver/app/agent/agent_api/internal/types"
	"beaver/app/agent/agent_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAgentModelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAgentModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentModelLogic {
	return &DeleteAgentModelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAgentModelLogic) DeleteAgentModel(req *types.DeleteAgentModelReq) (*types.DeleteAgentModelRes, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return nil, errors.New("user id required")
	}
	if strings.TrimSpace(req.ModelId) == "" {
		return nil, errors.New("modelId required")
	}

	result := l.svcCtx.DB.Where("model_id = ? AND user_id = ?", strings.TrimSpace(req.ModelId), strings.TrimSpace(req.UserID)).
		Delete(&agent_models.AgentUserModel{})
	if result.Error != nil {
		l.Errorf("DeleteAgentModel failed: %v", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("model not found")
	}

	return &types.DeleteAgentModelRes{}, nil
}
