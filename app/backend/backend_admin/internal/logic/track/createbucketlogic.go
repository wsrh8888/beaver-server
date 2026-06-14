package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBucketLogic {
	return &CreateBucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBucketLogic) CreateBucket(req *types.CreateBucketReq) (resp *types.CreateBucketRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.AdminCreateBucket(l.ctx, &platform_rpc.AdminCreateBucketReq{
		Name:        req.Name,
		Description: req.Description,
		CreateUser:  req.UserID,
	})
	if err != nil {
		l.Errorf("创建 Bucket 失败: %v", err)
		return nil, err
	}
	return &types.CreateBucketRes{BucketId: rpcRes.BucketId}, nil
}
