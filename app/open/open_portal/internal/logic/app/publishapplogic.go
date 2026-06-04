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

	// 查询应用
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限访问")
	}

	if app.EnableRobot == 1 {
		if err := ensurePortalAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, &app); err != nil {
			logx.Errorf("发布应用时 Robot 未就绪: app_id=%s err=%v", req.AppID, err)
			return nil, errors.New("发布失败：智能机器人未创建成功，请先开启 robot 能力后重试")
		}
	}

	// 更新应用状态为已发布
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
