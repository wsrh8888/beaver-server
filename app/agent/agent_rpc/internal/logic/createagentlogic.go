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

	"github.com/google/uuid"
)

type CreateAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentLogic {
	return &CreateAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("create_agent", ctx),
	}
}

func (l *CreateAgentLogic) CreateAgent(in *agent_rpc.CreateAgentReq) (*agent_rpc.CreateAgentRes, error) {
	if in == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(in.UserId) == "" {
		return nil, errors.New("user_id required")
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, errors.New("name required")
	}

	agent := agent_models.Agent{
		AgentID:     "agent_" + strings.ReplaceAll(uuid.NewString(), "-", ""),
		UserID:      in.UserId,
		Name:        name,
		Avatar:      strings.TrimSpace(in.Avatar),
		Description: strings.TrimSpace(in.Description),
		Status:      1,
	}
	if err := l.svcCtx.DB.WithContext(l.ctx).Create(&agent).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "创建智能体入库失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, err
	}

	return &agent_rpc.CreateAgentRes{
		AgentId:     agent.AgentID,
		Name:        agent.Name,
		Avatar:      agent.Avatar,
		Description: agent.Description,
	}, nil
}
