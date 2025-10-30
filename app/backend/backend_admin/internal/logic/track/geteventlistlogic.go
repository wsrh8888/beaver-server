package logic

import (
	"context"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/track/track_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取事件列表
func NewGetEventListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventListLogic {
	return &GetEventListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventListLogic) GetEventList(req *types.GetEventListReq) (resp *types.GetEventListRes, err error) {
	// 构建查询条件
	db := l.svcCtx.DB.Model(&track_models.TrackEvent{}).Preload("BucketModel")

	// Bucket Id 筛选
	if req.BucketID != "" {
		db = db.Where("bucket_id = ?", req.BucketID)
	}

	// 事件名称筛选
	if req.EventName != "" {
		db = db.Where("event_name LIKE ?", "%"+req.EventName+"%")
	}

	// 操作筛选
	if req.Action != "" {
		db = db.Where("action = ?", req.Action)
	}

	// 用户ID筛选
	if req.UserFilter != "" {
		db = db.Where("user_id = ?", req.UserFilter)
	}

	// 平台筛选
	if req.Platform != "" {
		db = db.Where("platform = ?", req.Platform)
	}

	// 时间范围筛选
	if req.StartTime > 0 {
		db = db.Where("timestamp >= ?", req.StartTime)
	}

	if req.EndTime > 0 {
		db = db.Where("timestamp <= ?", req.EndTime)
	}

	// 查询总数
	var total int64
	if err = db.Count(&total).Error; err != nil {
		logx.Errorf("查询事件总数失败: %v", err)
		return nil, err
	}

	// 分页查询
	var events []track_models.TrackEvent
	offset := (req.Page - 1) * req.PageSize
	if err = db.Offset(offset).Limit(req.PageSize).Order("timestamp DESC").Find(&events).Error; err != nil {
		logx.Errorf("查询事件列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	list := make([]types.GetEventListResItem, 0, len(events))
	for _, event := range events {
		bucketName := ""
		if event.BucketModel != nil {
			bucketName = event.BucketModel.Name
		}

		userID := ""
		if event.UserID != nil {
			userID = *event.UserID
		}

		list = append(list, types.GetEventListResItem{
			Id:         uint(event.Id),
			EventName:  event.EventName,
			Action:     event.Action,
			UserID:     userID,
			BucketID:   event.BucketID,
			BucketName: bucketName,
			Platform:   event.Platform,
			DeviceID:   event.DeviceID,
			Data:       string(event.Data),
			Timestamp:  event.Timestamp,
			CreatedAt:  time.Time(event.CreatedAt).Format(time.RFC3339),
		})
	}

	resp = &types.GetEventListRes{
		List:  list,
		Total: total,
	}

	return
}
