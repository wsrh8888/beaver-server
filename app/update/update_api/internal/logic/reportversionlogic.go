package logic

import (
	"context"
	"fmt"

	"beaver/app/update/update_api/internal/svc"
	"beaver/app/update/update_api/internal/types"
	"beaver/app/update/update_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ReportVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportVersionLogic {
	return &ReportVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportVersionLogic) ReportVersion(req *types.ReportVersionReq) (resp *types.ReportVersionRes, err error) {
	// H5架构特殊处理
	if req.ArchID == 0 {
		logx.Infof("H5版本上报: AppID=%s, Version=%s, DeviceID=%s", req.AppID, req.Version, req.DeviceID)
		return &types.ReportVersionRes{}, nil
	}

	// 1. 验证应用是否存在
	var app update_models.UpdateApp
	if err := l.svcCtx.DB.Where("app_id = ? AND is_active = ?", req.AppID, true).First(&app).Error; err != nil {
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

	// 3. 查找是否已有该设备的记录
	var existingReport update_models.UpdateReport
	err = l.svcCtx.DB.Where("device_id = ? AND app_id = ? AND architecture_id = ?",
		req.DeviceID, req.AppID, architecture.Id).First(&existingReport).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		logx.Errorf("查询设备记录失败: %v", err)
		return nil, fmt.Errorf("查询设备记录失败")
	}

	if err == gorm.ErrRecordNotFound {
		// 没有记录，创建新记录
		report := update_models.UpdateReport{
			UserID:         req.UserID,
			DeviceID:       req.DeviceID,
			AppID:          req.AppID,
			ArchitectureID: architecture.Id,
			Version:        req.Version,
		}

		if err := l.svcCtx.DB.Create(&report).Error; err != nil {
			logx.Errorf("创建版本上报记录失败: %v", err)
			return nil, fmt.Errorf("创建上报记录失败")
		}
	} else {
		// 有记录，更新现有记录
		updates := map[string]interface{}{
			"user_id": req.UserID,
			"version": req.Version,
		}

		if err := l.svcCtx.DB.Model(&existingReport).Updates(updates).Error; err != nil {
			logx.Errorf("更新版本上报记录失败: %v", err)
			return nil, fmt.Errorf("更新上报记录失败")
		}
	}

	logx.Infof("版本上报成功: AppID=%s, PlatformID=%d, ArchID=%d, Version=%s, DeviceID=%s",
		req.AppID, req.PlatformID, req.ArchID, req.Version, req.DeviceID)

	return &types.ReportVersionRes{}, nil
}
