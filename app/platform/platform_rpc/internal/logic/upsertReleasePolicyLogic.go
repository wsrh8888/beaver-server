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

type UpsertReleasePolicyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpsertReleasePolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpsertReleasePolicyLogic {
	return &UpsertReleasePolicyLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpsertReleasePolicyLogic) UpsertReleasePolicy(in *platform_rpc.UpsertReleasePolicyReq) (*platform_rpc.UpsertReleasePolicyRes, error) {
	if in.AppId == "" || in.ArchitectureId == 0 {
		return nil, status.Error(codes.InvalidArgument, "应用ID和架构ID不能为空")
	}
	if in.RolloutPercent > 100 {
		return nil, status.Error(codes.InvalidArgument, "灰度比例不能超过100")
	}
	if in.GrayVersionId > 0 && in.RolloutPercent == 0 {
		return nil, status.Error(codes.InvalidArgument, "配置灰度版本时必须设置灰度比例")
	}
	if in.StableVersionId == 0 {
		return nil, status.Error(codes.InvalidArgument, "请选择正式版")
	}

	if err := l.ensureVersionBelongs(in.ArchitectureId, in.StableVersionId); err != nil {
		return nil, err
	}
	if in.GrayVersionId > 0 {
		if err := l.ensureVersionBelongs(in.ArchitectureId, in.GrayVersionId); err != nil {
			return nil, err
		}
	}

	var policy platform_models.UpdateReleasePolicy
	err := l.svcCtx.DB.Where("architecture_id = ?", in.ArchitectureId).First(&policy).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		policy = platform_models.UpdateReleasePolicy{
			AppID:           in.AppId,
			ArchitectureID:  uint(in.ArchitectureId),
			StableVersionID: uint(in.StableVersionId),
			GrayVersionID:   uint(in.GrayVersionId),
			RolloutPercent:  uint(in.RolloutPercent),
			MinVersion:      in.MinVersion,
			ForceUpdate:     in.ForceUpdate,
			IsActive:        in.IsActive,
		}
		if err := l.svcCtx.DB.Create(&policy).Error; err != nil {
			return nil, status.Error(codes.Internal, "创建发版策略失败")
		}
		return &platform_rpc.UpsertReleasePolicyRes{Id: uint64(policy.Id)}, nil
	}
	if err != nil {
		return nil, err
	}

	policy.StableVersionID = uint(in.StableVersionId)
	policy.GrayVersionID = uint(in.GrayVersionId)
	policy.RolloutPercent = uint(in.RolloutPercent)
	policy.MinVersion = in.MinVersion
	policy.ForceUpdate = in.ForceUpdate
	policy.IsActive = in.IsActive
	if err := l.svcCtx.DB.Save(&policy).Error; err != nil {
		return nil, status.Error(codes.Internal, "更新发版策略失败")
	}
	return &platform_rpc.UpsertReleasePolicyRes{Id: uint64(policy.Id)}, nil
}

func (l *UpsertReleasePolicyLogic) ensureVersionBelongs(architectureID, versionID uint64) error {
	var version platform_models.UpdateVersion
	if err := l.svcCtx.DB.First(&version, versionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Error(codes.NotFound, "版本不存在")
		}
		return err
	}
	if version.ArchitectureID != uint(architectureID) {
		return status.Error(codes.InvalidArgument, "版本不属于当前架构")
	}
	return nil
}

func versionLabel(db *gorm.DB, versionID uint) string {
	if versionID == 0 {
		return ""
	}
	var v platform_models.UpdateVersion
	if err := db.First(&v, versionID).Error; err != nil {
		return ""
	}
	return v.Version
}
