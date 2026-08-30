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

type ListAgentModelsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAgentModelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAgentModelsLogic {
	return &ListAgentModelsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAgentModelsLogic) ListAgentModels(req *types.ListAgentModelsReq) (*types.ListAgentModelsRes, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return nil, errors.New("user id required")
	}

	var rows []agent_models.AgentUserModel
	err := l.svcCtx.DB.Where("user_id = ?", strings.TrimSpace(req.UserID)).
		Order("id desc").
		Find(&rows).Error
	if err != nil {
		l.Errorf("ListAgentModels failed: %v", err)
		return nil, err
	}

	list := make([]types.AgentModelItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toAgentModelItem(row))
	}

	return &types.ListAgentModelsRes{List: list}, nil
}
