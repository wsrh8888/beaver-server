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

type UpdateArchitectureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新架构信息
func NewUpdateArchitectureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArchitectureLogic {
	return &UpdateArchitectureLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateArchitectureLogic) UpdateArchitecture(req *types.UpdateArchitectureReq) (resp *types.UpdateArchitectureRes, err error) {
	// 检查架构是否存在
	var arch update_models.UpdateArchitecture
	if err := l.svcCtx.DB.First(&arch, req.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("architecture not found")
		}
		return nil, err
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.Description != "" {
		updates["description"] = req.Description
	}
	updates["is_active"] = req.IsActive

	// 执行更新
	if err := l.svcCtx.DB.Model(&arch).Updates(updates).Error; err != nil {
		logx.Errorf("Failed to update architecture: %v", err)
		return nil, err
	}

	return &types.UpdateArchitectureRes{}, nil
}
