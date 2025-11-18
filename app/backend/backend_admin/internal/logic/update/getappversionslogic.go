package logic

import (
	"beaver/app/update/update_models"
	"context"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppVersionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取应用下所有版本
func NewGetAppVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppVersionsLogic {
	return &GetAppVersionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppVersionsLogic) GetAppVersions(req *types.GetAppVersionsReq) (resp *types.GetAppVersionsRes, err error) {
	// 构建查询条件
	query := l.svcCtx.DB.Model(&update_models.UpdateArchitecture{})

	// 应用ID过滤
	if req.AppID != "" {
		query = query.Where("app_id = ?", req.AppID)
	}

	// 只查询活跃的架构
	query = query.Where("is_active = ?", true)

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logx.Errorf("Failed to count architectures: %v", err)
		return nil, fmt.Errorf("获取架构总数失败")
	}

	// 分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize

	// 查询架构列表
	var architectures []update_models.UpdateArchitecture
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&architectures).Error; err != nil {
		logx.Errorf("Failed to get architectures: %v", err)
		return nil, fmt.Errorf("获取架构列表失败")
	}

	// 转换为响应格式
	var architectureInfos []types.GetAppVersionsItem
	for _, arch := range architectures {
		// 查询该架构下的版本
		var versions []update_models.UpdateVersion
		if err := l.svcCtx.DB.Where("architecture_id = ?", arch.Id).Order("created_at DESC").Find(&versions).Error; err != nil {
			logx.Errorf("Failed to get versions for architecture %d: %v", arch.Id, err)
			continue
		}

		// 转换为简化版本信息
		var simpleVersions []types.GetAppVersionsVersionItem
		for _, version := range versions {
			simpleVersions = append(simpleVersions, types.GetAppVersionsVersionItem{
				VersionID: uint(version.Id),
				Version:   version.Version,
			})
		}

		architectureInfo := types.GetAppVersionsItem{
			ArchitectureID: uint(arch.Id),
			ArchID:         arch.ArchID,
			Description:    arch.Description,
			Versions:       simpleVersions,
		}
		architectureInfos = append(architectureInfos, architectureInfo)
	}

	return &types.GetAppVersionsRes{
		Total:         total,
		Architectures: architectureInfos,
	}, nil
}
