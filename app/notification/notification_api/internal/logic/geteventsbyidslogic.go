package logic

import (
	"context"
	"time"

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventsByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按ID拉取通知事件明细
func NewGetEventsByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventsByIdsLogic {
	return &GetEventsByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventsByIdsLogic) GetEventsByIds(req *types.GetEventsByIdsReq) (resp *types.GetEventsByIdsRes, err error) {
	resp = &types.GetEventsByIdsRes{Events: []types.EventItem{}}

	if len(req.EventIDs) == 0 {
		return resp, nil
	}

	var rows []notification_models.NotificationEvent
	if err = l.svcCtx.DB.WithContext(l.ctx).
		Where("event_id IN ?", req.EventIDs).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, ev := range rows {
		resp.Events = append(resp.Events, types.EventItem{
			EventID:    ev.EventID,
			EventType:  ev.EventType,
			Category:   ev.Category,
			Version:    ev.Version,
			FromUserID: derefString(ev.FromUserID),
			TargetID:   derefString(ev.TargetID),
			TargetType: ev.TargetType,
			Payload:    string(ev.Payload),
			Priority:   int32(ev.Priority),
			Status:     int32(ev.Status),
			DedupHash:  ev.DedupHash,
			CreatedAt:  time.Time(ev.CreatedAt).UnixMilli(),
			UpdatedAt:  time.Time(ev.UpdatedAt).UnixMilli(),
		})
	}

	return resp, nil
}

func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
