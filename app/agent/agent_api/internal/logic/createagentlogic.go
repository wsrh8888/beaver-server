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

type CreateAgentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentLogic {
	return &CreateAgentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAgentLogic) CreateAgent(req *types.CreateAgentReq) (*types.CreateAgentRes, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return nil, errors.New("user id required")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("name required")
	}

	res, err := l.svcCtx.AgentRpc.CreateAgent(l.ctx, &agent_rpc.CreateAgentReq{
		UserId:      req.UserID,
		Name:        strings.TrimSpace(req.Name),
		Avatar:      strings.TrimSpace(req.Avatar),
		Description: strings.TrimSpace(req.Description),
	})
	if err != nil {
		l.Errorf("CreateAgent rpc failed: %v", err)
		return nil, err
	}

	return &types.CreateAgentRes{
		Agent: types.AgentItem{
			AgentId:     res.AgentId,
			Name:        res.Name,
			Avatar:      res.Avatar,
			Description: res.Description,
			Status:      1,
		},
	}, nil
}
