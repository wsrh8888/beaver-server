package track_public

import (
	"context"
	"encoding/json"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/app/platform/platform_models"

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

func (l *ReportEventsLogic) ReportEvents(req *types.ReportEventsReq) (*types.ReportEventsRes, error) {
	events := make([]platform_models.TrackEvent, 0, len(req.Events))

	for _, event := range req.Events {
		eventData := map[string]interface{}{
			"platform":  event.Platform,
			"timestamp": event.Timestamp,
		}
		if event.Platform == "" {
			delete(eventData, "platform")
		}
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
			l.Logger.Errorf("marshal event data failed: %v", err)
			continue
		}

		var userID *string
		if req.UserID != "" {
			userID = &req.UserID
		}

		events = append(events, platform_models.TrackEvent{
			EventName: event.EventName,
			Action:    event.Action,
			UserID:    userID,
			DeviceID:  event.DeviceID,
			BucketID:  event.BucketID,
			Platform:  event.Platform,
			Timestamp: event.Timestamp,
			Data:      datatypes.JSON(jsonData),
		})
	}

	go func(eventsToSave []platform_models.TrackEvent) {
		err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			return tx.CreateInBatches(eventsToSave, 100).Error
		})
		if err != nil {
			l.Logger.Errorf("save events failed: %v", err)
		}
	}(events)

	return &types.ReportEventsRes{}, nil
}
