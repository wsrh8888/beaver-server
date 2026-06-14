package logic

import (
	"context"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	uuid_util "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCreateBucketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminCreateBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateBucketLogic {
	return &AdminCreateBucketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminCreateBucketLogic) AdminCreateBucket(in *platform_rpc.AdminCreateBucketReq) (*platform_rpc.AdminCreateBucketRes, error) {
	bucketID := uuid_util.NewV4().String()
	bucket := &platform_models.TrackBucket{
		BucketID:    bucketID,
		Name:        in.Name,
		Description: in.Description,
		CreateUser:  in.CreateUser,
		IsActive:    true,
	}
	if err := l.svcCtx.DB.Create(bucket).Error; err != nil {
		l.Errorf("创建 Bucket 失败: %v", err)
		return nil, err
	}
	return &platform_rpc.AdminCreateBucketRes{BucketId: bucketID}, nil
}
