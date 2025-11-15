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

type AddVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加新版本
func NewAddVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddVersionLogic {
	return &AddVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddVersionLogic) AddVersion(req *types.AddVersionReq) (resp *types.AddVersionRes, err error) {
	// 检查架构是否存在
	var arch update_models.UpdateArchitecture
	if err := l.svcCtx.DB.First(&arch, req.ArchitectureID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("architecture not found")
		}
		return nil, err
	}

	// 解析发布时间
	if err != nil {
		return nil, fmt.Errorf("invalid release date format")
	}

	// 创建新版本
	version := update_models.UpdateVersion{
		ArchitectureID: req.ArchitectureID,
		Version:        req.Version,
		FileKey:        req.FileName,
		Description:    req.Description,
		ReleaseNotes:   req.ReleaseNotes,
	}

	if err := l.svcCtx.DB.Create(&version).Error; err != nil {
		logx.Errorf("Failed to create version: %v", err)
		return nil, err
	}

	return &types.AddVersionRes{
		VersionID: uint(version.Id),
	}, nil
}
