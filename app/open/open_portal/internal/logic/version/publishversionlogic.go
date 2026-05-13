package version

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发布版本
func NewPublishVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishVersionLogic {
	return &PublishVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishVersionLogic) PublishVersion(req *types.PublishVersionReq) (resp *types.PublishVersionRes, err error) {
	// 解析 versionId
	var versionID uint
	fmt.Sscanf(req.VersionID, "%d", &versionID)

	// 查询版本
	var version open_models.OpenAppVersion
	if err := l.svcCtx.DB.First(&version, versionID).Error; err != nil {
		return nil, err
	}

	// 验证应用所有权
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", version.AppID).First(&app).Error; err != nil {
		return nil, err
	}

	// 验证 UserID 是否为应用所有者
	if app.OwnerUserID != req.UserID {
		return nil, errors.New("无权操作此应用")
	}

	// 只有审核通过的状态才能发布
	if version.Status != "approved" {
		return nil, fmt.Errorf("当前状态(%s)不允许发布，需要先审核通过", version.Status)
	}

	// 更新版本状态为已发布
	if err := l.svcCtx.DB.Model(&version).Update("status", "published").Error; err != nil {
		return nil, err
	}

	// 更新应用状态为已发布
	if err := l.svcCtx.DB.Model(&open_models.OpenApp{}).Where("app_id = ?", version.AppID).Update("status", 1).Error; err != nil {
		return nil, err
	}

	return &types.PublishVersionRes{}, nil
}
