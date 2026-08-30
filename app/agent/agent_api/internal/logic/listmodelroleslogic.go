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

type ListModelRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListModelRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListModelRolesLogic {
	return &ListModelRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListModelRolesLogic) ListModelRoles(req *types.ListModelRolesReq) (*types.ListModelRolesRes, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return nil, errors.New("user id required")
	}

	var rows []agent_models.AgentModelRole
	err := l.svcCtx.DB.Where("status = ?", 1).Find(&rows).Error
	if err != nil {
		l.Errorf("ListModelRoles failed: %v", err)
		return nil, err
	}

	list := make([]types.AgentModelRoleItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, types.AgentModelRoleItem{
			Role:            row.Role,
			ModelId:         row.ModelID,
			FallbackModelId: row.FallbackModelID,
			Status:          int32(row.Status),
		})
	}

	return &types.ListModelRolesRes{List: list}, nil
}
