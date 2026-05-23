package app

import (
	"context"
	"errors"

	"beaver/app/open/constants"
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

	// 创建 defaultScopes map 提高查找效率
	defaultMap := make(map[string]bool)
	for _, s := range constants.DefaultScopes {
		defaultMap[string(s)] = true
	}

	scopes := make([]types.ScopeInfo, 0, len(constants.AllScopes))
	for _, scope := range constants.AllScopes {
		scopeStr := string(scope)
		scopes = append(scopes, types.ScopeInfo{
			Scope:       scopeStr,
			Name:        scopeStr,
			Description: constants.ScopeDescription[scope],
			Required:    defaultMap[scopeStr],
		})
	}

	return &types.GetAppScopesRes{
		Scopes: scopes,
	}, nil
}
