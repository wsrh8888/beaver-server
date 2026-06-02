package app

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新应用
func NewUpdateAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppLogic {
	return &UpdateAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAppLogic) UpdateApp(req *types.UpdateAppReq) (resp *types.UpdateAppRes, err error) {
	if _, err := l.svcCtx.RequireDeveloper(req.UserID); err != nil {
		return nil, err
	}

	// 查询应用并验证权限
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限")
	}

	// 3. 构建更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	// 4. 执行更新
	if len(updates) == 0 {
		return &types.UpdateAppRes{}, nil
	}

	if err := l.svcCtx.DB.Model(&app).Updates(updates).Error; err != nil {
		logx.Errorf("更新应用失败: %v", err)
		return nil, errors.New("更新失败")
	}

	logx.Infof("应用更新成功: app_id=%s", req.AppID)

	return &types.UpdateAppRes{}, nil
}
