package bot

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBotConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新 Bot 配置
func NewUpdateBotConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBotConfigLogic {
	return &UpdateBotConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBotConfigLogic) UpdateBotConfig(req *types.UpdateBotConfigReq) (resp *types.UpdateBotConfigRes, err error) {
	// 1. 检查应用是否存在且属于当前用户
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	// 2. 查询或创建 Bot 配置
	var botConfig open_models.OpenBotConfig
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&botConfig).Error; err != nil {
		// 不存在则创建
		botConfig = open_models.OpenBotConfig{
			AppID:            req.AppID,
			BotName:          "",
			BotAvatar:        "",
			BotDescription:   "",
			UsageGuide:       "",
			EnableSingleChat: 1,
			EnableGroupChat:  1,
			EnableAtMention:  1,
			EnableMenu:       0,
			MenuItems:        "[]",
			AutoReplyRules:   "[]",
			Commands:         "[]",
			Status:           1,
		}
	}

	// 3. 更新字段（只更新传入的字段）
	if req.BotName != "" {
		botConfig.BotName = req.BotName
	}
	if req.BotAvatar != "" {
		botConfig.BotAvatar = req.BotAvatar
	}
	if req.BotDescription != "" {
		botConfig.BotDescription = req.BotDescription
	}
	if req.UsageGuide != "" {
		botConfig.UsageGuide = req.UsageGuide
	}
	if req.EnableSingleChat != nil {
		if *req.EnableSingleChat {
			botConfig.EnableSingleChat = 1
		} else {
			botConfig.EnableSingleChat = 0
		}
	}
	if req.EnableGroupChat != nil {
		if *req.EnableGroupChat {
			botConfig.EnableGroupChat = 1
		} else {
			botConfig.EnableGroupChat = 0
		}
	}
	if req.EnableAtMention != nil {
		if *req.EnableAtMention {
			botConfig.EnableAtMention = 1
		} else {
			botConfig.EnableAtMention = 0
		}
	}
	if req.EnableMenu != nil {
		if *req.EnableMenu {
			botConfig.EnableMenu = 1
		} else {
			botConfig.EnableMenu = 0
		}
	}
	if req.MenuItems != "" {
		botConfig.MenuItems = req.MenuItems
	}
	if req.AutoReplyRules != nil {
		rulesJSON, _ := json.Marshal(req.AutoReplyRules)
		botConfig.AutoReplyRules = string(rulesJSON)
	}
	if req.Commands != nil {
		commandsJSON, _ := json.Marshal(req.Commands)
		botConfig.Commands = string(commandsJSON)
	}
	if req.Status != nil {
		botConfig.Status = *req.Status
	}

	// 4. 保存配置
	if botConfig.ID == 0 {
		// 新建
		if err := l.svcCtx.DB.Create(&botConfig).Error; err != nil {
			return nil, errors.New("创建 Bot 配置失败")
		}
	} else {
		// 更新
		if err := l.svcCtx.DB.Save(&botConfig).Error; err != nil {
			return nil, errors.New("更新 Bot 配置失败")
		}
	}

	return &types.UpdateBotConfigRes{}, nil
}
