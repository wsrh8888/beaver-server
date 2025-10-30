package logic

import (
	"context"

	"beaver/app/ai/ai_models"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAIBotsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取AI机器人列表
func NewGetAIBotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAIBotsLogic {
	return &GetAIBotsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAIBotsLogic) GetAIBots(req *types.GetAIBotsReq) (resp *types.GetAIBotsRes, err error) {
	// 默认分页参数
	page := 1
	limit := 20
	if req.Page > 0 {
		page = req.Page
	}
	if req.Limit > 0 {
		limit = req.Limit
	}

	// 查询AI机器人列表
	var bots []ai_models.AIModelConfig
	var total int64
	err = l.svcCtx.DB.Model(&ai_models.AIModelConfig{}).
		Order("id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&bots).Error
	if err != nil {
		return nil, err
	}

	// 获取总数
	err = l.svcCtx.DB.Model(&ai_models.AIModelConfig{}).Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 转换为API响应格式
	var botConfigs []types.AIBotConfig
	for _, bot := range bots {
		botConfigs = append(botConfigs, types.AIBotConfig{
			ID:          bot.Id,
			Name:        bot.Name,
			Provider:    bot.Provider,
			FileName:    bot.FileName,
			Description: bot.Description,
			MaxTokens:   bot.MaxTokens,
			Temperature: bot.Temperature,
			IsEnabled:   bot.IsEnabled,
			SystemMsg:   bot.SystemMsg,
			WelcomeMsg:  bot.WelcomeMsg,
			Features:    bot.Features,
		})
	}

	return &types.GetAIBotsRes{
		Total: total,
		Bots:  botConfigs,
	}, nil
}
