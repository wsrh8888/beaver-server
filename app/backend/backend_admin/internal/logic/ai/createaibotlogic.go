package logic

import (
	"context"
	"strconv"

	"beaver/app/ai/ai_models"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAIBotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建AI机器人
func NewCreateAIBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAIBotLogic {
	return &CreateAIBotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAIBotLogic) CreateAIBot(req *types.CreateAIBotReq) (resp *types.CreateAIBotRes, err error) {
	// 创建AI机器人配置
	bot := ai_models.AIModelConfig{
		Name:        req.Bot.Name,
		Provider:    req.Bot.Provider,
		FileName:    req.Bot.FileName,
		Description: req.Bot.Description,
		MaxTokens:   req.Bot.MaxTokens,
		Temperature: req.Bot.Temperature,
		IsEnabled:   req.Bot.IsEnabled,
		SystemMsg:   req.Bot.SystemMsg,
		WelcomeMsg:  req.Bot.WelcomeMsg,
		Features:    req.Bot.Features,
	}

	// 保存到数据库
	err = l.svcCtx.DB.Create(&bot).Error
	if err != nil {
		return nil, err
	}

	return &types.CreateAIBotRes{
		BotID: strconv.FormatUint(uint64(bot.Id), 10),
	}, nil
}
