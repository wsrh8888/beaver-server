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
	if _, err := l.svcCtx.RequireDeveloper(req.UserID); err != nil {
		return nil, err
	}

	// 查询应用
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限访问")
	}

	// 2. TODO: OpenBotModel 已重构为群机器人模型，应用维度的 Bot 配置功能暂时禁用
	// 原逻辑：检查并创建默认 Bot 配置
	logx.Infof("应用发布前检查（Bot 配置功能待实现）: app_id=%s", req.AppID)

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
