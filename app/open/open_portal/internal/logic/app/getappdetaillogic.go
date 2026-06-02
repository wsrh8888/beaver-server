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
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}

	// 查询应用详情
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限访问")
	}

	// 2. 对 AppSecret 进行掩码处理（只显示前8位和后8位）
	maskedSecret := maskSecret(app.AppSecret)

	return &types.GetAppDetailRes{
		App: types.AppInfo{
			AppID:       app.AppID,
			Name:        app.Name,
			Description: app.Description,
			Icon:        app.Icon,
			AppSecret:   maskedSecret,
			Status:      app.Status,
			// 能力开关
			EnableRobot:   app.EnableRobot,
			EnableOAuth:   app.EnableOAuth,
			EnableWebhook: app.EnableWebhook,
			CreatedAt:     app.CreatedAt.Unix(),
		},
	}, nil
}

// maskSecret 对密钥进行掩码处理
func maskSecret(secret string) string {
	if len(secret) <= 16 {
		return "****"
	}
	return secret[:8] + "****" + secret[len(secret)-8:]
}
