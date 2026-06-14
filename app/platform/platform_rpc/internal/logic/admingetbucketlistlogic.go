package logic

import (
	"context"
	"time"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetBucketListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminGetBucketListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetBucketListLogic {
	return &AdminGetBucketListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminGetBucketListLogic) AdminGetBucketList(in *platform_rpc.AdminGetBucketListReq) (*platform_rpc.AdminGetBucketListRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	db := l.svcCtx.DB.Model(&platform_models.TrackBucket{})
	if in.Keyword != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+in.Keyword+"%", "%"+in.Keyword+"%")
	}
	if in.IsActive != nil {
		db = db.Where("is_active = ?", *in.IsActive)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("查询 Bucket 总数失败: %v", err)
		return nil, err
	}

	var buckets []platform_models.TrackBucket
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&buckets).Error; err != nil {
		l.Errorf("查询 Bucket 列表失败: %v", err)
		return nil, err
	}

	list := make([]*platform_rpc.AdminBucketItem, 0, len(buckets))
	for _, bucket := range buckets {
		list = append(list, &platform_rpc.AdminBucketItem{
			BucketId:    bucket.BucketID,
			Name:        bucket.Name,
			Description: bucket.Description,
			CreateUser:  bucket.CreateUser,
			IsActive:    bucket.IsActive,
			CreatedAt:   time.Time(bucket.CreatedAt).Format(time.RFC3339),
			UpdatedAt:   time.Time(bucket.UpdatedAt).Format(time.RFC3339),
		})
	}

	return &platform_rpc.AdminGetBucketListRes{List: list, Total: total}, nil
}
