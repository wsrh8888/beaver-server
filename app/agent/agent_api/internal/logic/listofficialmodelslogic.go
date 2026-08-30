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

type ListOfficialModelsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOfficialModelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOfficialModelsLogic {
	return &ListOfficialModelsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOfficialModelsLogic) ListOfficialModels(req *types.ListOfficialModelsReq) (*types.ListOfficialModelsRes, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return nil, errors.New("user id required")
	}

	var rows []agent_models.AgentOfficialModel
	err := l.svcCtx.DB.Where("status = ?", 1).
		Order("sort asc, id asc").
		Find(&rows).Error
	if err != nil {
		l.Errorf("ListOfficialModels failed: %v", err)
		return nil, err
	}

	// 写死虚拟 Auto，始终置顶；库里若误插同 id 则跳过
	list := make([]types.AgentOfficialModelItem, 0, len(rows)+1)
	list = append(list, types.AgentOfficialModelItem{
		ModelId:   agent_models.OfficialModelIDAuto,
		Name:      "Auto",
		Provider:  "beaver",
		ModelName: agent_models.OfficialModelIDAuto,
		Endpoint:  "",
		Tier:      agent_models.ModelTierBalanced,
		Capabilities: types.AgentModelCapabilities{
			Vision:           true,
			Tools:            true,
			Reasoning:        true,
			StructuredOutput: true,
			LongContext:      true,
		},
		Sort:   -1,
		Status: 1,
	})

	for _, row := range rows {
		if row.ModelID == agent_models.OfficialModelIDAuto {
			continue
		}
		caps := row.Capabilities
		list = append(list, types.AgentOfficialModelItem{
			ModelId:   row.ModelID,
			Name:      row.Name,
			Provider:  row.Provider,
			ModelName: row.ModelName,
			Endpoint:  row.Endpoint,
			Tier:      row.Tier,
			Capabilities: types.AgentModelCapabilities{
				Vision:           caps.Vision,
				Tools:            caps.Tools,
				Reasoning:        caps.Reasoning,
				StructuredOutput: caps.StructuredOutput,
				LongContext:      caps.LongContext,
			},
			Sort:   row.Sort,
			Status: int32(row.Status),
		})
	}

	return &types.ListOfficialModelsRes{List: list}, nil
}
