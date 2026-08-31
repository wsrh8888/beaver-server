package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/agent/agent_models"
	"beaver/app/agent/agent_rpc/internal/svc"
	"beaver/app/agent/agent_rpc/types/agent_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type GetAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentLogic {
	return &GetAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_agent", ctx),
	}
}

func (l *GetAgentLogic) GetAgent(in *agent_rpc.GetAgentReq) (*agent_rpc.GetAgentRes, error) {
	if in == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(in.UserId) == "" {
		return nil, errors.New("user_id required")
	}
	if strings.TrimSpace(in.AgentId) == "" {
		return nil, errors.New("agent_id required")
	}

	var agent agent_models.Agent
	err := l.svcCtx.DB.WithContext(l.ctx).
		Where("agent_id = ? AND user_id = ?", in.AgentId, in.UserId).
		First(&agent).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("agent not found")
	}
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询智能体失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, err
	}

	return &agent_rpc.GetAgentRes{
		AgentId:     agent.AgentID,
		Name:        agent.Name,
		Avatar:      agent.Avatar,
		Description: agent.Description,
		Status:      int32(agent.Status),
	}, nil
}
