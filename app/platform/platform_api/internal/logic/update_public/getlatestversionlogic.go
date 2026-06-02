package update_public

import (
	"context"
	"fmt"
	"strings"

	"beaver/app/file/file_rpc/types/file_rpc"
	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/app/platform/platform_models"

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
		logx.Infof("H5获取最新版本: AppID=%s, Version=%s, DeviceID=%s", req.AppID, req.Version, req.DeviceID)
		return &types.GetLatestVersionRes{HasUpdate: false}, nil
	}

	var app platform_models.UpdateApp
	if err := l.svcCtx.DB.Where("app_id = ? AND is_active = ?", req.AppID, true).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("应用不存在或已停用")
		}
		logx.Errorf("查询应用失败: %v", err)
		return nil, fmt.Errorf("查询应用失败")
	}

	var architecture platform_models.UpdateArchitecture
	if err := l.svcCtx.DB.Where("app_id = ? AND platform_id = ? AND arch_id = ? AND is_active = ?",
		req.AppID, req.PlatformID, req.ArchID, true).First(&architecture).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("不支持的平台或架构")
		}
		logx.Errorf("查询架构失败: %v", err)
		return nil, fmt.Errorf("查询架构失败")
	}

	var latestVersion platform_models.UpdateVersion
	if err := l.svcCtx.DB.Where("architecture_id = ?", architecture.Id).
		Order("created_at DESC").First(&latestVersion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &types.GetLatestVersionRes{HasUpdate: false}, nil
		}
		logx.Errorf("查询最新版本失败: %v", err)
		return nil, fmt.Errorf("查询最新版本失败")
	}

	if !l.isNewerVersion(latestVersion.Version, req.Version) {
		return &types.GetLatestVersionRes{HasUpdate: false}, nil
	}

	resp := &types.GetLatestVersionRes{
		HasUpdate:      true,
		ArchitectureID: uint(architecture.Id),
		Version:        latestVersion.Version,
		FileKey:        latestVersion.FileKey,
		Description:    latestVersion.Description,
		ReleaseNotes:   latestVersion.ReleaseNotes,
	}

	fileDetail, err := l.svcCtx.FileRpc.GetFileDetail(l.ctx, &file_rpc.GetFileDetailReq{
		FileKey: latestVersion.FileKey,
	})
	if err != nil {
		logx.Errorf("获取文件详情失败: %v", err)
	} else {
		resp.Size = fileDetail.Size
		resp.MD5 = fileDetail.Md5
	}

	logx.Infof("获取最新版本成功: AppID=%s, PlatformID=%d, ArchID=%d, CurrentVersion=%s, LatestVersion=%s",
		req.AppID, req.PlatformID, req.ArchID, req.Version, latestVersion.Version)

	return resp, nil
}

func (l *GetLatestVersionLogic) isNewerVersion(latestVersion, currentVersion string) bool {
	return strings.Compare(latestVersion, currentVersion) > 0
}
