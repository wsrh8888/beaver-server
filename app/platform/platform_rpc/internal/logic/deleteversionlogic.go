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

type DeleteVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVersionLogic {
	return &DeleteVersionLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteVersionLogic) DeleteVersion(in *platform_rpc.DeleteVersionReq) (*platform_rpc.DeleteVersionRes, error) {
	if in.VersionId == 0 {
		return nil, status.Error(codes.InvalidArgument, "版本ID不能为空")
	}

	var version platform_models.UpdateVersion
	if err := l.svcCtx.DB.First(&version, in.VersionId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "版本不存在")
		}
		return nil, err
	}

	var policies []platform_models.UpdateReleasePolicy
	if err := l.svcCtx.DB.Find(&policies).Error; err != nil {
		return nil, err
	}
	vid := uint(in.VersionId)
	for _, p := range policies {
		if p.StableVersionID == vid || p.GrayVersionID == vid {
			return nil, status.Error(codes.FailedPrecondition, "版本已被发版策略引用，请先调整策略")
		}
	}

	if err := l.svcCtx.DB.Delete(&version).Error; err != nil {
		l.Errorf("删除版本失败: %v", err)
		return nil, status.Error(codes.Internal, "删除版本失败")
	}

	return &platform_rpc.DeleteVersionRes{}, nil
}
