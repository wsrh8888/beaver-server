package app

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleAppCapabilityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 启用/禁用应用能力（对标飞书）
func NewToggleAppCapabilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleAppCapabilityLogic {
	return &ToggleAppCapabilityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleAppCapabilityLogic) ToggleAppCapability(req *types.ToggleAppCapabilityReq) (resp *types.ToggleAppCapabilityRes, err error) {
	// 1. 查询应用
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	// 2. 根据能力类型更新对应的开关
	var enabled bool
	switch req.Capability {
	case "bot":
		if req.Enable {
			app.EnableBot = 1
			enabled = true
			// 如果启用 Bot 能力，检查是否存在 Bot 配置，不存在则创建默认配置
			var botConfig open_models.OpenBotModel
			if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&botConfig).Error; err != nil {
				// 创建默认 Bot 配置
				botConfig = open_models.OpenBotModel{
					AppID:            req.AppID,
					Name:        app.Name,
					Avatar:      app.Icon,
					Description: app.Description,
					EnableSingleChat: 1,
					EnableGroupChat:  1,
					EnableAtMention:  1,
					Status:           1,
				}
				if err := l.svcCtx.DB.Create(&botConfig).Error; err != nil {
					logx.Errorf("创建默认 Bot 配置失败: %v", err)
					return nil, errors.New("创建 Bot 配置失败")
				}
				logx.Infof("为应用 %s 创建默认 Bot 配置", req.AppID)
			}
		} else {
			app.EnableBot = 0
			enabled = false
		}
	case "oauth":
		if req.Enable {
			app.EnableOAuth = 1
			enabled = true
		} else {
			app.EnableOAuth = 0
			enabled = false
		}
	case "webhook":
		if req.Enable {
			app.EnableWebhook = 1
			enabled = true
		} else {
			app.EnableWebhook = 0
			enabled = false
		}
	default:
		return nil, errors.New("不支持的能力类型")
	}

	// 3. 保存更新
	if err := l.svcCtx.DB.Save(&app).Error; err != nil {
		logx.Errorf("更新应用能力失败: %v", err)
		return nil, errors.New("更新应用能力失败")
	}

	logx.Infof("应用 %s 的 %s 能力已%s", req.AppID, req.Capability, map[bool]string{true: "启用", false: "禁用"}[req.Enable])

	return &types.ToggleAppCapabilityRes{
		Enabled: enabled,
	}, nil
}
