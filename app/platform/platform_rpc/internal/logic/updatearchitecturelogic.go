package logic

import (
	"context"
	"errors"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UpdateArchitectureLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateArchitectureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArchitectureLogic {
	return &UpdateArchitectureLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateArchitectureLogic) UpdateArchitecture(in *platform_rpc.UpdateArchitectureReq) (*platform_rpc.UpdateArchitectureRes, error) {
	var arch platform_models.UpdateArchitecture
	if err := l.svcCtx.DB.First(&arch, in.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "架构不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{"is_active": in.IsActive}
	if in.Description != "" {
		updates["description"] = in.Description
	}
	if err := l.svcCtx.DB.Model(&arch).Updates(updates).Error; err != nil {
		l.Errorf("更新架构失败: %v", err)
		return nil, status.Error(codes.Internal, "更新架构失败")
	}

	return &platform_rpc.UpdateArchitectureRes{}, nil
}
