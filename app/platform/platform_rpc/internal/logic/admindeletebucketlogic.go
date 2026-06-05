package logic

import (
	"context"
	"errors"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteBucketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminDeleteBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteBucketLogic {
	return &AdminDeleteBucketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminDeleteBucketLogic) AdminDeleteBucket(in *platform_rpc.AdminDeleteBucketReq) (*platform_rpc.AdminDeleteBucketRes, error) {
	result := l.svcCtx.DB.Where("bucket_id = ?", in.BucketId).Delete(&platform_models.TrackBucket{})
	if result.Error != nil {
		l.Errorf("删除 Bucket 失败: %v", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("Bucket 不存在")
	}
	return &platform_rpc.AdminDeleteBucketRes{}, nil
}
