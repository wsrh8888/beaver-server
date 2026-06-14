package logic

import (
	"context"
	"errors"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateBucketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminUpdateBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateBucketLogic {
	return &AdminUpdateBucketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminUpdateBucketLogic) AdminUpdateBucket(in *platform_rpc.AdminUpdateBucketReq) (*platform_rpc.AdminUpdateBucketRes, error) {
	updates := map[string]interface{}{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Description != "" {
		updates["description"] = in.Description
	}
	if in.IsActive != nil {
		updates["is_active"] = *in.IsActive
	}
	if len(updates) == 0 {
		return &platform_rpc.AdminUpdateBucketRes{}, nil
	}

	result := l.svcCtx.DB.Model(&platform_models.TrackBucket{}).Where("bucket_id = ?", in.BucketId).Updates(updates)
	if result.Error != nil {
		l.Errorf("更新 Bucket 失败: %v", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("Bucket 不存在")
	}
	return &platform_rpc.AdminUpdateBucketRes{}, nil
}
