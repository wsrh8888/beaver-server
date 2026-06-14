package update_public

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/app/platform/platform_models"
	"beaver/core/coreupdate"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetLatestVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLatestVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLatestVersionLogic {
	return &GetLatestVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLatestVersionLogic) GetLatestVersion(req *types.GetLatestVersionReq) (*types.GetLatestVersionRes, error) {
	if req.ArchID == 0 {
		return &types.GetLatestVersionRes{HasUpdate: false}, nil
	}

	var app platform_models.UpdateApp
	if err := l.svcCtx.DB.Where("app_id = ? AND is_active = ?", req.AppID, true).First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("应用不存在或已停用")
		}
		return nil, fmt.Errorf("查询应用失败")
	}

	var architecture platform_models.UpdateArchitecture
	if err := l.svcCtx.DB.Where("app_id = ? AND platform_id = ? AND arch_id = ? AND is_active = ?",
		req.AppID, req.PlatformID, req.ArchID, true).First(&architecture).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("不支持的平台或架构")
		}
		return nil, fmt.Errorf("查询架构失败")
	}

	stableID, grayID, rollout, minVer, forceFlag, active := l.loadPolicy(architecture)

	resolved := coreupdate.Resolve(coreupdate.ResolveInput{
		AppID:           req.AppID,
		ArchitectureID:  architecture.Id,
		DeviceID:        req.DeviceID,
		UserID:          req.UserID,
		CurrentVersion:  req.Version,
		StableVersionID: stableID,
		GrayVersionID:   grayID,
		RolloutPercent:  rollout,
		MinVersion:      minVer,
		ForceUpdate:     forceFlag,
		PolicyActive:    active,
	})

	if resolved.TargetVersionID == 0 {
		return &types.GetLatestVersionRes{HasUpdate: false}, nil
	}

	var target platform_models.UpdateVersion
	if err := l.svcCtx.DB.First(&target, resolved.TargetVersionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.GetLatestVersionRes{HasUpdate: false}, nil
		}
		return nil, fmt.Errorf("查询目标版本失败")
	}

	if !coreupdate.CompareVersion(target.Version, req.Version) {
		return &types.GetLatestVersionRes{HasUpdate: false}, nil
	}

	resp := &types.GetLatestVersionRes{
		HasUpdate:      true,
		ForceUpdate:    resolved.ForceUpdate,
		ArchitectureID: architecture.Id,
		Version:        target.Version,
		FileUrl:        target.FileUrl,
		Description:    target.Description,
		ReleaseNotes:   target.ReleaseNotes,
	}

	logx.Infof("检查更新: app=%s arch=%d current=%s target=%s gray=%v force=%v",
		req.AppID, architecture.Id, req.Version, target.Version, resolved.InGrayRollout, resolved.ForceUpdate)

	return resp, nil
}

func (l *GetLatestVersionLogic) loadPolicy(arch platform_models.UpdateArchitecture) (stableID, grayID, rollout uint, minVer string, forceFlag, active bool) {
	var policy platform_models.UpdateReleasePolicy
	err := l.svcCtx.DB.Where("architecture_id = ? AND is_active = ?", arch.Id, true).First(&policy).Error
	if err == nil {
		return policy.StableVersionID, policy.GrayVersionID, policy.RolloutPercent, policy.MinVersion, policy.ForceUpdate, true
	}

	var latest platform_models.UpdateVersion
	if err := l.svcCtx.DB.Where("architecture_id = ?", arch.Id).Order("created_at DESC").First(&latest).Error; err == nil {
		return latest.Id, 0, 0, "", false, true
	}
	return 0, 0, 0, "", false, false
}
