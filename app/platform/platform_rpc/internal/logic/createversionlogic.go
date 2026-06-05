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

type CreateVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateVersionLogic {
	return &CreateVersionLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateVersionLogic) CreateVersion(in *platform_rpc.CreateVersionReq) (*platform_rpc.CreateVersionRes, error) {
	var arch platform_models.UpdateArchitecture
	if err := l.svcCtx.DB.First(&arch, in.ArchitectureId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "架构不存在")
		}
		return nil, err
	}

	version := platform_models.UpdateVersion{
		ArchitectureID: uint(in.ArchitectureId),
		Version:        in.Version,
		FileKey:        in.FileKey,
		Description:    in.Description,
		ReleaseNotes:   in.ReleaseNotes,
	}
	if err := l.svcCtx.DB.Create(&version).Error; err != nil {
		l.Errorf("创建版本失败: %v", err)
		return nil, status.Error(codes.Internal, "创建版本失败")
	}

	return &platform_rpc.CreateVersionRes{VersionId: uint64(version.Id)}, nil
}
