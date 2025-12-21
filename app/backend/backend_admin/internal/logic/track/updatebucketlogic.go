package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/track/track_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新Bucket
func NewUpdateBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBucketLogic {
	return &UpdateBucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBucketLogic) UpdateBucket(req *types.UpdateBucketReq) (resp *types.UpdateBucketRes, err error) {
	// 构建更新字段
	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}

	if req.Description != "" {
		updates["description"] = req.Description
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// 执行更新操作
	result := l.svcCtx.DB.Model(&track_models.TrackBucket{}).Where("bucket_id = ?", req.BucketId).Updates(updates)

	if result.Error != nil {
		logx.Errorf("更新Bucket失败: %v", result.Error)
		return nil, result.Error
	}

	// 检查是否更新了记录
	if result.RowsAffected == 0 {
		logx.Errorf("未找到要更新的Bucket: %s", req.BucketId)
		return nil, result.Error
	}

	resp = &types.UpdateBucketRes{}

	return
}
