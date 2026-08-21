package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/agent/agent_models"
	"beaver/app/agent/agent_rpc/internal/svc"
	"beaver/app/agent/agent_rpc/types/agent_rpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentLogic {
	return &CreateAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
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
		l.Errorf("CreateAgent insert failed: %v", err)
		return nil, err
	}

	return &agent_rpc.CreateAgentRes{
		AgentId:     agent.AgentID,
		Name:        agent.Name,
		Avatar:      agent.Avatar,
		Description: agent.Description,
	}, nil
}
