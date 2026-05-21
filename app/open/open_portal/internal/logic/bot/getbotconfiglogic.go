package bot

import (
	"context"
	"encoding/json"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBotConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 Bot 配置
func NewGetBotConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBotConfigLogic {
	return &GetBotConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBotConfigLogic) GetBotConfig(req *types.GetBotConfigReq) (resp *types.GetBotConfigRes, err error) {
	// 1. 查询 Bot 配置
	var botConfig open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&botConfig).Error; err != nil {
		// 如果不存在，返回默认配置
		return &types.GetBotConfigRes{
			Config: types.BotConfigInfo{
				AppID:            req.AppID,
				BotName:          "",
				BotAvatar:        "",
				BotDescription:   "",
				UsageGuide:       "",
				EnableSingleChat: true,
				EnableGroupChat:  true,
				EnableAtMention:  true,
				EnableMenu:       false,
				MenuItems:        "[]",
				AutoReplyRules:   []string{},
				Commands:         []string{},
				Status:           1,
			},
		}, nil
	}

	// 2. 解析 JSON 字段
	var autoReplyRules []string
	var commands []string
	if botConfig.AutoReplyRules != "" {
		json.Unmarshal([]byte(botConfig.AutoReplyRules), &autoReplyRules)
	}
	if botConfig.Commands != "" {
		json.Unmarshal([]byte(botConfig.Commands), &commands)
	}

	// 3. 返回配置
	return &types.GetBotConfigRes{
		Config: types.BotConfigInfo{
			AppID:            botConfig.AppID,
			BotName:          botConfig.Name,
			BotAvatar:        botConfig.Avatar,
			BotDescription:   botConfig.Description,
			UsageGuide:       botConfig.UsageGuide,
			EnableSingleChat: botConfig.EnableSingleChat == 1,
			EnableGroupChat:  botConfig.EnableGroupChat == 1,
			EnableAtMention:  botConfig.EnableAtMention == 1,
			EnableMenu:       botConfig.EnableMenu == 1,
			MenuItems:        botConfig.MenuItems,
			AutoReplyRules:   autoReplyRules,
			Commands:         commands,
			Status:           botConfig.Status,
		},
	}, nil
}
