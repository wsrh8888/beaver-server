package logic

import (
	"context"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/track/track_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBucketListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Bucket列表
func NewGetBucketListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBucketListLogic {
	return &GetBucketListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBucketListLogic) GetBucketList(req *types.GetBucketListReq) (resp *types.GetBucketListRes, err error) {
	// 构建查询条件
	db := l.svcCtx.DB.Model(&track_models.TrackBucket{})

	// 关键词搜索 (名称或描述)
	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 状态筛选
	if req.IsActive != nil {
		db = db.Where("is_active = ?", *req.IsActive)
	}

	// 查询总数
	var total int64
	if err = db.Count(&total).Error; err != nil {
		logx.Errorf("查询Bucket总数失败: %v", err)
		return nil, err
	}

	// 分页查询
	var buckets []track_models.TrackBucket
	offset := (req.Page - 1) * req.PageSize
	if err = db.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&buckets).Error; err != nil {
		logx.Errorf("查询Bucket列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	list := make([]types.GetBucketListItem, 0, len(buckets))
	for _, bucket := range buckets {
		list = append(list, types.GetBucketListItem{
			BucketId:    bucket.BucketID,
			Name:        bucket.Name,
			Description: bucket.Description,
			CreateUser:  bucket.CreateUser,
			IsActive:    bucket.IsActive,
			CreatedAt:   time.Time(bucket.CreatedAt).Format(time.RFC3339),
			UpdatedAt:   time.Time(bucket.UpdatedAt).Format(time.RFC3339),
		})
	}

	resp = &types.GetBucketListRes{
		List:  list,
		Total: total,
	}

	return
}
