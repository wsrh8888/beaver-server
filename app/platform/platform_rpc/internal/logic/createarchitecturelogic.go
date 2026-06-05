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

type CreateArchitectureLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateArchitectureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateArchitectureLogic {
	return &CreateArchitectureLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateArchitectureLogic) CreateArchitecture(in *platform_rpc.CreateArchitectureReq) (*platform_rpc.CreateArchitectureRes, error) {
	var app platform_models.UpdateApp
	if err := l.svcCtx.DB.Where("app_id = ?", in.AppId).First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "应用不存在")
		}
		return nil, err
	}

	var existing platform_models.UpdateArchitecture
	if err := l.svcCtx.DB.Where("app_id = ? AND platform_id = ? AND arch_id = ?", in.AppId, in.PlatformId, in.ArchId).First(&existing).Error; err == nil {
		return nil, status.Error(codes.AlreadyExists, "架构已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	arch := platform_models.UpdateArchitecture{
		AppID:       in.AppId,
		PlatformID:  uint(in.PlatformId),
		ArchID:      uint(in.ArchId),
		Description: in.Description,
		IsActive:    true,
	}
	if err := l.svcCtx.DB.Create(&arch).Error; err != nil {
		l.Errorf("创建架构失败: %v", err)
		return nil, status.Error(codes.Internal, "创建架构失败")
	}

	return &platform_rpc.CreateArchitectureRes{Id: uint64(arch.Id)}, nil
}
