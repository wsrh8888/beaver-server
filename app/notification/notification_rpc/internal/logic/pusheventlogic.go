package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/notification/notification_models"
	"beaver/app/notification/notification_rpc/internal/svc"
	"beaver/app/notification/notification_rpc/types/notification_rpc"
	mqwsconst "beaver/common/const/rocketmq"
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
	fmt.Println("1111111111111111111111111111")
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
	go func(etcdAddr string, toUsers []string, eventID string, eventVersion int64) {
		// 构建表更新数据 - 包含版本号，让前端知道具体同步哪些数据
		var tableUpdates []map[string]interface{}

		// 通知事件表更新
		eventUpdates := map[string]interface{}{
			"table": "notification_event",
			"data": []map[string]interface{}{
				{
					"version": eventVersion,
					"eventId": eventID,
				},
			},
		}
		tableUpdates = append(tableUpdates, eventUpdates)

		// 为每个收件人推送表更新通知
		for _, uid := range toUsers {
			// 为每个用户生成独立的收件箱版本号用于通知
			inboxVersion := l.svcCtx.VersionGen.GetNextVersion(notification_models.VersionScopeInboxPerUser, "user_id", uid)

			// 通知收件箱表更新
			inboxUpdates := map[string]interface{}{
				"table":  "notification_inbox",
				"userId": uid,
				"data": []map[string]interface{}{
					{
						"version": inboxVersion,
						"eventId": eventID,
					},
				},
			}
			// 为每个用户创建独立的tableUpdates副本
			userTableUpdates := append([]map[string]interface{}{}, tableUpdates...)
			userTableUpdates = append(userTableUpdates, inboxUpdates)

			payload := map[string]interface{}{
				"command":  wsCommandConst.NOTIFICATION,
				"type":     wsTypeConst.NotificationReceive,
				"senderId": in.FromUserId,
				"targetId": uid,
				"body": map[string]interface{}{
					"tableUpdates": userTableUpdates,
				},
				"conversationId": "",
			}
			l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload)
		}
	}(l.svcCtx.Config.Etcd.Hosts[0], in.ToUserIds, eventID, eventVersion)

	return &notification_rpc.PushEventRes{
		EventId: eventID,
		Version: eventVersion,
	}, nil
}
