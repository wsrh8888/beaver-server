package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"beaver/app/track/track_api/internal/svc"
	"beaver/app/track/track_api/internal/types"
	"beaver/app/track/track_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReportEventsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportEventsLogic {
	return &ReportEventsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ReportEvents 处理上报埋点事件请求
func (l *ReportEventsLogic) ReportEvents(req *types.ReportEventsReq) (resp *types.ReportEventsRes, err error) {
	// 准备批量插入的事件数据
	events := make([]track_models.TrackEvent, 0, len(req.Events))
	fmt.Println("111111111111111111111111", req.UserID)

	for _, event := range req.Events {
		// 构建事件数据JSON
		eventData := map[string]interface{}{
			"platform":  event.Platform,
			"timestamp": event.Timestamp,
		}

		// 移除空值
		if event.Platform == "" {
			delete(eventData, "platform")
		}

		// 如果有额外数据，解析并合并
		if event.Data != "" {
			var extraData map[string]interface{}
			if err := json.Unmarshal([]byte(event.Data), &extraData); err == nil {
				for k, v := range extraData {
					eventData[k] = v
				}
			}
		}

		jsonData, err := json.Marshal(eventData)
		if err != nil {
			l.Logger.Errorf("Failed to marshal event data: %v", err)
			continue
		}

		// 处理UserID指针
		var userID *string
		if req.UserID != "" {
			userID = &req.UserID
		}

		// 创建事件记录
		trackEvent := track_models.TrackEvent{
			EventName: event.EventName,
			Action:    event.Action,
			UserID:    userID,
			DeviceID:  event.DeviceID,
			BucketID:  event.BucketID,
			Platform:  event.Platform,
			Timestamp: event.Timestamp,
			Data:      datatypes.JSON(jsonData),
		}

		events = append(events, trackEvent)
	}

	// 异步批量保存事件
	go func(eventsToSave []track_models.TrackEvent) {
		err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			result := tx.CreateInBatches(eventsToSave, 100)
			if result.Error != nil {
				l.Logger.Errorf("Failed to save events: %v", result.Error)
				return result.Error
			}
			return nil
		})

		if err != nil {
			l.Logger.Errorf("Transaction failed when saving events: %v", err)
		}
	}(events)

	return &types.ReportEventsRes{}, nil
}
