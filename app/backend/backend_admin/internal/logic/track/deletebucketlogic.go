package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBucketLogic {
	return &DeleteBucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBucketLogic) DeleteBucket(req *types.DeleteBucketReq) (resp *types.DeleteBucketRes, err error) {
	_, err = l.svcCtx.PlatformRpc.AdminDeleteBucket(l.ctx, &platform_rpc.AdminDeleteBucketReq{
		BucketId: req.BucketId,
	})
	if err != nil {
		l.Errorf("删除 Bucket 失败: %v", err)
		return nil, err
	}
	return &types.DeleteBucketRes{}, nil
}
