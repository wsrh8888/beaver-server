package version

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建新版本
func NewCreateVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateVersionLogic {
	return &CreateVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateVersionLogic) CreateVersion(req *types.CreateVersionReq) (resp *types.CreateVersionRes, err error) {
	// 验证应用所有权
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		return nil, err
	}

	// 验证 UserID 是否为应用所有者
	if app.OwnerUserID != req.UserID {
		return nil, errors.New("无权操作此应用")
	}

	// 根据应用当前启用的能力自动填充 capabilities
	var capabilities []string
	if app.EnableBot == 1 {
		capabilities = append(capabilities, "bot")
	}
	if app.EnableOAuth == 1 {
		capabilities = append(capabilities, "oauth")
	}
	if app.EnableWebhook == 1 {
		capabilities = append(capabilities, "webhook")
	}

	capabilitiesJSON := "[]"
	if len(capabilities) > 0 {
		data, _ := json.Marshal(capabilities)
		capabilitiesJSON = string(data)
	}

	// 创建版本记录
	version := open_models.OpenAppVersion{
		AppID:        req.AppID,
		Version:      req.Version,
		Description:  req.Description,
		Visibility:   req.Visibility,
		Status:       "draft",
		Capabilities: capabilitiesJSON,
		CreatedBy:    req.UserID,
	}

	if err := l.svcCtx.DB.Create(&version).Error; err != nil {
		return nil, err
	}

	return &types.CreateVersionRes{
		VersionID: fmt.Sprintf("%d", version.ID),
	}, nil
}
