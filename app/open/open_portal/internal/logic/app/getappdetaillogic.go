package app

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取应用详情
func NewGetAppDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppDetailLogic {
	return &GetAppDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppDetailLogic) GetAppDetail(req *types.GetAppDetailReq) (resp *types.GetAppDetailRes, err error) {
	// 1. 从 header 获取当前用户 ID
	userID := l.ctx.Value("userId")
	if userID == nil {
		return nil, errors.New("未登录")
	}

	// 2. 查询应用详情
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, userID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限访问")
	}

	return &types.GetAppDetailRes{
		App: types.AppInfo{
			AppID:       app.AppID,
			Name:        app.Name,
			Description: app.Description,
			Icon:        app.Icon,
			Status:      app.Status,
			WebhookURL:  app.WebhookURL,
			CreatedAt:   app.CreatedAt.Unix(),
		},
	}, nil
}
