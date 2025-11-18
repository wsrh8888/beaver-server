package logic

import (
	"beaver/app/update/update_models"
	"context"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AddArchitectureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加新架构
func NewAddArchitectureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddArchitectureLogic {
	return &AddArchitectureLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddArchitectureLogic) AddArchitecture(req *types.AddArchitectureReq) (resp *types.AddArchitectureRes, err error) {
	// 检查应用是否存在
	var app update_models.UpdateApp
	if err := l.svcCtx.DB.Where("uuid = ?", req.AppID).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("app not found")
		}
		return nil, err
	}
	fmt.Println("1111111111111111111111111111111")
	// 检查是否已存在相同的架构
	var existingArch update_models.UpdateArchitecture
	if err := l.svcCtx.DB.Where("app_id = ? AND platform_id = ? AND arch_id = ?",
		req.AppID, req.PlatformID, req.ArchID).First(&existingArch).Error; err == nil {
		return nil, fmt.Errorf("architecture already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	fmt.Println("22222222222222222222222222222222222")

	// 创建新架构
	arch := update_models.UpdateArchitecture{
		AppID:       req.AppID,
		PlatformID:  req.PlatformID,
		ArchID:      req.ArchID,
		Description: req.Description,
		IsActive:    true, // 默认为活跃状态
	}
	fmt.Println("333333333333333333333333333333333333333")

	if err := l.svcCtx.DB.Create(&arch).Error; err != nil {
		logx.Errorf("Failed to create architecture: %v", err)
		return nil, err
	}
	fmt.Println("44444444444444444444444444444444444")

	return &types.AddArchitectureRes{
		Id: uint(arch.Id),
	}, nil
}
