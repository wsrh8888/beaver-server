package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/track/track_models"
	uuid_util "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建Bucket
func NewCreateBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBucketLogic {
	return &CreateBucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBucketLogic) CreateBucket(req *types.CreateBucketReq) (resp *types.CreateBucketRes, err error) {
	// 生成唯一 Bucket ID
	bucketID := uuid_util.NewV4().String()

	// 创建 Bucket 记录
	bucket := &track_models.TrackBucket{
		BucketID:    bucketID,
		Name:        req.Name,
		Description: req.Description,
		CreateUser:  req.UserID,
		IsActive:    true, // 默认激活状态
	}

	// 插入数据库
	if err = l.svcCtx.DB.Create(bucket).Error; err != nil {
		logx.Errorf("创建Bucket失败: %v", err)
		return nil, err
	}

	resp = &types.CreateBucketRes{
		BucketId: bucketID,
	}

	return
}
