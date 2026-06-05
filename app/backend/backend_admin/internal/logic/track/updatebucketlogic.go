package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBucketLogic {
	return &UpdateBucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBucketLogic) UpdateBucket(req *types.UpdateBucketReq) (resp *types.UpdateBucketRes, err error) {
	_, err = l.svcCtx.PlatformRpc.AdminUpdateBucket(l.ctx, &platform_rpc.AdminUpdateBucketReq{
		BucketId:    req.BucketId,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	})
	if err != nil {
		l.Errorf("更新 Bucket 失败: %v", err)
		return nil, err
	}
	return &types.UpdateBucketRes{}, nil
}
