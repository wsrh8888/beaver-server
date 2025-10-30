package logic

import (
	"context"
	"strconv"

	"beaver/app/ai/ai_models"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAIBotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新AI机器人
func NewUpdateAIBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAIBotLogic {
	return &UpdateAIBotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAIBotLogic) UpdateAIBot(req *types.UpdateAIBotReq) (resp *types.UpdateAIBotRes, err error) {
	// 解析机器人ID
	botID, err := strconv.ParseUint(strconv.FormatUint(uint64(req.Bot.ID), 10), 10, 64)
	if err != nil {
		return nil, err
	}

	// 查询机器人是否存在
	var bot ai_models.AIModelConfig
	err = l.svcCtx.DB.First(&bot, botID).Error
	if err != nil {
		return nil, err
	}

	// 更新机器人信息
	bot.Name = req.Bot.Name
	bot.Provider = req.Bot.Provider
	bot.FileName = req.Bot.FileName
	bot.Description = req.Bot.Description
	bot.MaxTokens = req.Bot.MaxTokens
	bot.Temperature = req.Bot.Temperature
	bot.IsEnabled = req.Bot.IsEnabled
	bot.SystemMsg = req.Bot.SystemMsg
	bot.WelcomeMsg = req.Bot.WelcomeMsg
	bot.Features = req.Bot.Features

	// 保存更新
	err = l.svcCtx.DB.Save(&bot).Error
	if err != nil {
		return nil, err
	}

	return &types.UpdateAIBotRes{}, nil
}
