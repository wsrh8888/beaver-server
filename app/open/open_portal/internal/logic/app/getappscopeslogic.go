package app

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppScopesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取应用权限列表
func NewGetAppScopesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppScopesLogic {
	return &GetAppScopesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppScopesLogic) GetAppScopes(req *types.GetAppScopesReq) (resp *types.GetAppScopesRes, err error) {
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

	// 3. 解析已授权的权限列表
	var enabledScopes []string
	if app.Scopes != "" {
		json.Unmarshal([]byte(app.Scopes), &enabledScopes)
	}

	// 4. 构建所有权限列表（标记是否已启用）
	scopes := make([]types.ScopeInfo, 0)
	for _, scope := range open_models.AllScopes {
		enabled := false
		for _, s := range enabledScopes {
			if string(scope) == s {
				enabled = true
				break
			}
		}

		// 判断是否是必需权限（默认权限）
		required := false
		for _, s := range open_models.DefaultScopes {
			if scope == s {
				required = true
				break
			}
		}

		scopes = append(scopes, types.ScopeInfo{
			Scope:       string(scope),
			Name:        string(scope),
			Description: open_models.ScopeDescription[scope],
			Enabled:     enabled,
			Required:    required,
		})
	}

	return &types.GetAppScopesRes{
		Scopes: scopes,
	}, nil
}
