package logic

import (
	"context"
	"time"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetEventListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminGetEventListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetEventListLogic {
	return &AdminGetEventListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminGetEventListLogic) AdminGetEventList(in *platform_rpc.AdminGetEventListReq) (*platform_rpc.AdminGetEventListRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}

	db := l.svcCtx.DB.Model(&platform_models.TrackEvent{}).Preload("BucketModel")
	if in.BucketId != "" {
		db = db.Where("bucket_id = ?", in.BucketId)
	}
	if in.EventName != "" {
		db = db.Where("event_name LIKE ?", "%"+in.EventName+"%")
	}
	if in.Action != "" {
		db = db.Where("action = ?", in.Action)
	}
	if in.UserFilter != "" {
		db = db.Where("user_id = ?", in.UserFilter)
	}
	if in.Platform != "" {
		db = db.Where("platform = ?", in.Platform)
	}
	if in.StartTime > 0 {
		db = db.Where("timestamp >= ?", in.StartTime)
	}
	if in.EndTime > 0 {
		db = db.Where("timestamp <= ?", in.EndTime)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("查询事件总数失败: %v", err)
		return nil, err
	}

	var events []platform_models.TrackEvent
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("timestamp DESC").Find(&events).Error; err != nil {
		l.Errorf("查询事件列表失败: %v", err)
		return nil, err
	}

	list := make([]*platform_rpc.AdminEventItem, 0, len(events))
	for _, event := range events {
		bucketName := ""
		if event.BucketModel != nil {
			bucketName = event.BucketModel.Name
		}
		userID := ""
		if event.UserID != nil {
			userID = *event.UserID
		}
		dataStr := ""
		if event.Data != nil {
			dataBytes, _ := event.Data.MarshalJSON()
			dataStr = string(dataBytes)
		}
		list = append(list, &platform_rpc.AdminEventItem{
			Id:         uint64(event.Id),
			EventName:  event.EventName,
			Action:     event.Action,
			UserId:     userID,
			BucketId:   event.BucketID,
			BucketName: bucketName,
			Platform:   event.Platform,
			DeviceId:   event.DeviceID,
			Data:       dataStr,
			Timestamp:  event.Timestamp,
			CreatedAt:  time.Time(event.CreatedAt).Format(time.RFC3339),
		})
	}

	return &platform_rpc.AdminGetEventListRes{List: list, Total: total}, nil
}
