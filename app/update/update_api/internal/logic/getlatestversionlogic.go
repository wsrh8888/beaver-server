package logic

import (
	"context"
	"fmt"
	"strings"

	"beaver/app/file/file_rpc/types/file_rpc"
	"beaver/app/update/update_api/internal/svc"
	"beaver/app/update/update_api/internal/types"
	"beaver/app/update/update_models"

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

func (l *GetLatestVersionLogic) GetLatestVersion(req *types.GetLatestVersionReq) (resp *types.GetLatestVersionRes, err error) {
	fmt.Println("req.CityName:", req.CityName)
	// H5架构特殊处理
	if req.ArchID == 0 {
		logx.Infof("H5获取最新版本: AppID=%s, Version=%s, DeviceID=%s", req.AppID, req.Version, req.DeviceID)
		return &types.GetLatestVersionRes{
			HasUpdate: false,
		}, nil
	}

	// 1. 验证应用是否存在
	var app update_models.UpdateApp
	if err := l.svcCtx.DB.Where("uuid = ? AND is_active = ?", req.AppID, true).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("应用不存在或已停用")
		}
		logx.Errorf("查询应用失败: %v", err)
		return nil, fmt.Errorf("查询应用失败")
	}

	// 2. 验证架构是否存在
	var architecture update_models.UpdateArchitecture
	if err := l.svcCtx.DB.Where("app_id = ? AND platform_id = ? AND arch_id = ? AND is_active = ?",
		req.AppID, req.PlatformID, req.ArchID, true).First(&architecture).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("不支持的平台或架构")
		}
		logx.Errorf("查询架构失败: %v", err)
		return nil, fmt.Errorf("查询架构失败")
	}

	// 3. 查找该架构下的最新版本
	var latestVersion update_models.UpdateVersion
	if err := l.svcCtx.DB.Where("architecture_id = ?", architecture.Id).
		Order("created_at DESC").First(&latestVersion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 没有找到版本，返回无更新
			return &types.GetLatestVersionRes{
				HasUpdate: false,
			}, nil
		}
		logx.Errorf("查询最新版本失败: %v", err)
		return nil, fmt.Errorf("查询最新版本失败")
	}

	// 4. 比较版本号（简单字符串比较，实际项目中可能需要更复杂的版本比较逻辑）
	if !l.isNewerVersion(latestVersion.Version, req.Version) {
		// 当前版本已经是最新或更新
		return &types.GetLatestVersionRes{
			HasUpdate: false,
		}, nil
	}

	// 5. 通过file_rpc获取文件详情
	fileDetail, err := l.svcCtx.FileRpc.GetFileDetail(l.ctx, &file_rpc.GetFileDetailReq{
		FileKey: latestVersion.FileKey,
	})
	if err != nil {
		logx.Errorf("获取文件详情失败: %v", err)
		// 即使获取文件详情失败，也返回版本信息，只是没有文件大小和MD5
		resp = &types.GetLatestVersionRes{
			HasUpdate:      true,
			ArchitectureID: uint(architecture.Id),
			Version:        latestVersion.Version,
			FileKey:        latestVersion.FileKey,
			Size:           fileDetail.Size,
			MD5:            fileDetail.Md5,
			Description:    latestVersion.Description,
			ReleaseNotes:   latestVersion.ReleaseNotes,
		}
	} else {
		// 6. 构建响应
		resp = &types.GetLatestVersionRes{
			HasUpdate:      true,
			ArchitectureID: uint(architecture.Id),
			Version:        latestVersion.Version,
			FileKey:        latestVersion.FileKey,
			Size:           fileDetail.Size,
			MD5:            fileDetail.Md5,
			Description:    latestVersion.Description,
			ReleaseNotes:   latestVersion.ReleaseNotes,
		}
	}

	logx.Infof("获取最新版本成功: AppID=%s, PlatformID=%d, ArchID=%d, CurrentVersion=%s, LatestVersion=%s",
		req.AppID, req.PlatformID, req.ArchID, req.Version, latestVersion.Version)

	return resp, nil
}

// 简单的版本比较逻辑（实际项目中可能需要更复杂的版本比较）
func (l *GetLatestVersionLogic) isNewerVersion(latestVersion, currentVersion string) bool {
	// 这里使用简单的字符串比较，实际项目中应该使用语义化版本比较
	return strings.Compare(latestVersion, currentVersion) > 0
}
