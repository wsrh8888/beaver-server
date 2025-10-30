package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/track/track_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除Bucket
func NewDeleteBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBucketLogic {
	return &DeleteBucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBucketLogic) DeleteBucket(req *types.DeleteBucketReq) (resp *types.DeleteBucketRes, err error) {
	// 执行软删除
	result := l.svcCtx.DB.Where("uuid = ?", req.UUID).Delete(&track_models.TrackBucket{})

	if result.Error != nil {
		logx.Errorf("删除Bucket失败: %v", result.Error)
		return nil, result.Error
	}

	// 检查是否删除了记录
	if result.RowsAffected == 0 {
		logx.Errorf("未找到要删除的Bucket: %s", req.UUID)
		return nil, result.Error
	}

	resp = &types.DeleteBucketRes{}

	return
}
