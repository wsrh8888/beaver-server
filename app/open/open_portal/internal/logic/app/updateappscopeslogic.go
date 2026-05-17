package app

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/open/constants"
	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppScopesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新应用权限
func NewUpdateAppScopesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppScopesLogic {
	return &UpdateAppScopesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAppScopesLogic) UpdateAppScopes(req *types.UpdateAppScopesReq) (resp *types.UpdateAppScopesRes, err error) {
	// 1. 从 header 获取当前用户 ID
	userID := l.ctx.Value("userId")
	if userID == nil {
		return nil, errors.New("未登录")
	}

	// 2. 查询应用
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, userID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限")
	}

	// 3. 验证权限合法性
	validScopes := make(map[string]bool)
	for _, scope := range constants.AllScopes {
		validScopes[string(scope)] = true
	}

	for _, scope := range req.Scopes {
		if !validScopes[scope] {
			return nil, errors.New("无效的权限: " + scope)
		}
	}

	// 4. 确保默认权限始终存在
	// 创建 req.Scopes 的 set
	scopeSet := make(map[string]bool)
	for _, s := range req.Scopes {
		scopeSet[s] = true
	}

	// 添加缺失的默认权限
	for _, s := range constants.DefaultScopes {
		scopeStr := string(s)
		if !scopeSet[scopeStr] {
			req.Scopes = append(req.Scopes, scopeStr)
		}
	}

	// 5. 序列化并保存
	scopesJSON, _ := json.Marshal(req.Scopes)
	if err := l.svcCtx.DB.Model(&app).Update("scopes", string(scopesJSON)).Error; err != nil {
		logx.Errorf("更新权限失败: %v", err)
		return nil, errors.New("更新权限失败")
	}

	logx.Infof("应用权限更新成功: app_id=%s, scopes=%v", req.AppID, req.Scopes)

	return &types.UpdateAppScopesRes{}, nil
}
