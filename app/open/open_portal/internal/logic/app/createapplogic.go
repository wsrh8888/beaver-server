package app

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建应用
func NewCreateAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAppLogic {
	return &CreateAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAppLogic) CreateApp(req *types.CreateAppReq) (resp *types.CreateAppRes, err error) {
	if _, err := l.svcCtx.RequireDeveloper(req.UserID); err != nil {
		return nil, err
	}

	// 生成 AppID 和 AppSecret
	appID := fmt.Sprintf("app_%s", uuid.New().String()[:8])
	appSecret := uuid.New().String() + uuid.New().String()

	// 3. 创建应用记录（状态为草稿，需要发布后才能被用户搜索到）
	app := open_models.OpenApp{
		AppID:       appID,
		AppSecret:   appSecret,
		Name:        req.Name,
		Description: req.Description,
		OwnerUserID: req.UserID,
		Status:      0, // 0=草稿，1=已发布，2=禁用
		// Icon:        req.Icon, // TODO: 数据库添加 icon 字段后启用
	}

	if err := l.svcCtx.DB.Create(&app).Error; err != nil {
		logx.Errorf("创建应用失败: %v", err)
		return nil, errors.New("创建应用失败")
	}

	logx.Infof("应用创建成功: app_id=%s, user_id=%s", appID, req.UserID)

	return &types.CreateAppRes{
		AppID:     appID,
		AppSecret: appSecret,
	}, nil
}
