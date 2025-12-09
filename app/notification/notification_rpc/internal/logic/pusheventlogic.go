package logic

import (
	"context"
	"errors"

	"beaver/app/notification/notification_models"
	"beaver/app/notification/notification_rpc/internal/svc"
	"beaver/app/notification/notification_rpc/types/notification_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PushEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushEventLogic {
	return &PushEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushEventLogic) PushEvent(in *notification_rpc.PushEventReq) (*notification_rpc.PushEventRes, error) {
	if len(in.ToUserIds) == 0 {
		return nil, errors.New("to_user_ids 不能为空")
	}
	if in.EventType == "" || in.Category == "" {
		return nil, errors.New("event_type 和 category 不能为空")
	}

	eventID := uuid.NewString()

	// 生成事件版本（全局）
	eventVersion := l.svcCtx.VersionGen.GetNextVersion(notification_models.VersionScopeEventGlobal, "", "")
	if eventVersion == -1 {
		return nil, errors.New("生成事件版本失败")
	}

	err := l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		ev := &notification_models.NotificationEvent{
			EventID:    eventID,
			EventType:  in.EventType,
			Category:   in.Category,
			Version:    eventVersion,
			TargetType: in.TargetType,
			Payload:    datatypes.JSON([]byte(in.PayloadJson)),
			Priority:   5,
			Status:     1,
			DedupHash:  in.DedupHash,
		}
		if in.FromUserId != "" {
			ev.FromUserID = &in.FromUserId
		}
		if in.TargetId != "" {
			ev.TargetID = &in.TargetId
		}

		if err := tx.Create(ev).Error; err != nil {
			return err
		}

		inboxRows := make([]notification_models.NotificationInbox, 0, len(in.ToUserIds))
		for _, uid := range in.ToUserIds {
			// 针对用户的收件箱版本（按 user_id 递增）
			inboxVersion := l.svcCtx.VersionGen.GetNextVersion(notification_models.VersionScopeInboxPerUser, "user_id", uid)
			if inboxVersion == -1 {
				return errors.New("生成收件箱版本失败")
			}

			inboxRows = append(inboxRows, notification_models.NotificationInbox{
				UserID:    uid,
				EventID:   eventID,
				EventType: in.EventType,
				Category:  in.Category,
				Version:   inboxVersion,
				IsRead:    false,
				ReadAt:    nil,
				Status:    1,
				Silent:    false,
			})
		}
		if len(inboxRows) > 0 {
			if err := tx.Create(&inboxRows).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// 异步通过 WS 通知收件人有新通知
	go func(etcdAddr string, toUsers []string, eventID string) {
		payload := map[string]interface{}{
			"eventId": eventID,
			"hasNew":  true,
		}
		for _, uid := range toUsers {
			ajax.SendMessageToWs(etcdAddr, wsCommandConst.NOTIFICATION, wsTypeConst.NotificationReceive, in.FromUserId, uid, payload, "")
		}
	}(l.svcCtx.Config.Etcd, in.ToUserIds, eventID)

	return &notification_rpc.PushEventRes{
		EventId: eventID,
		Version: eventVersion,
	}, nil
}
