package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/agent/agent_api/internal/svc"
	"beaver/app/agent/agent_api/internal/types"
	"beaver/app/agent/agent_rpc/types/agent_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentLogic {
	return &GetAgentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAgentLogic) GetAgent(req *types.GetAgentReq) (*types.GetAgentRes, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return nil, errors.New("user id required")
	}
	if strings.TrimSpace(req.AgentId) == "" {
		return nil, errors.New("agentId required")
	}

	res, err := l.svcCtx.AgentRpc.GetAgent(l.ctx, &agent_rpc.GetAgentReq{
		UserId:  req.UserID,
		AgentId: req.AgentId,
	})
	if err != nil {
		l.Errorf("GetAgent rpc failed: %v", err)
		return nil, err
	}

	return &types.GetAgentRes{
		Agent: types.AgentItem{
			AgentId:     res.AgentId,
			Name:        res.Name,
			Avatar:      res.Avatar,
			Description: res.Description,
			Status:      res.Status,
		},
	}, nil
}
