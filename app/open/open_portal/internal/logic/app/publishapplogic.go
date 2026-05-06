package app

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发布应用（对标飞书版本发布）
func NewPublishAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishAppLogic {
	return &PublishAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishAppLogic) PublishApp(req *types.PublishAppReq) (resp *types.PublishAppRes, err error) {
	// 1. 查询应用
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限访问")
	}

	// 2. 检查是否已配置 Bot（如果启用了 Bot 能力）
	var botConfig open_models.OpenBotConfig
	err = l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&botConfig).Error
	if err != nil {
		// 如果没有 Bot 配置，创建默认配置
		botConfig = open_models.OpenBotConfig{
			AppID:            req.AppID,
			BotName:          app.Name,
			BotAvatar:        app.Icon,
			BotDescription:   app.Description,
			EnableSingleChat: 1, // int 类型：1是 0否
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

	// 3. 更新应用状态为已发布
	if err := l.svcCtx.DB.Model(&app).Update("status", 1).Error; err != nil {
		logx.Errorf("发布应用失败: %v", err)
		return nil, errors.New("发布应用失败")
	}

	// 4. 记录版本信息（可选，后续可以创建 open_app_versions 表）
	logx.Infof("应用发布成功: app_id=%s, version=%s, user_id=%s", req.AppID, req.Version, req.UserID)

	return &types.PublishAppRes{
		Status: 1, // 已发布
	}, nil
}
