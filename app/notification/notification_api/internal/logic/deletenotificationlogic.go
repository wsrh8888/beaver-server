package logic

import (
	"context"

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type DeleteNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// 按事件ID删除单个通知
func NewDeleteNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotificationLogic {
	return &DeleteNotificationLogic{
		ctx:    ctx,
		logger: logger.New("delete_notification"),
		svcCtx: svcCtx,
	}
}

func (l *DeleteNotificationLogic) DeleteNotification(req *types.DeleteNotificationReq) (resp *types.DeleteNotificationRes, err error) {
	userId := req.UserID
	eventId := req.EventID

	// 将指定用户的指定通知标记为已删除
	result := l.svcCtx.DB.Model(&notification_models.NotificationInbox{}).
		Where("user_id = ? AND event_id = ?", userId, eventId).
		Update("is_deleted", true)

	if result.Error != nil {
		logx.WithContext(l.ctx).Errorf("删除通知失败: %v", result.Error)
		return nil, result.Error
	}

	resp = &types.DeleteNotificationRes{
		Success: result.RowsAffected > 0,
	}

	if result.RowsAffected > 0 {
		l.logger.Info(model.LogMsg{
			Text: "通知删除成功",
			Data: map[string]interface{}{
				"userId":  userId,
				"eventId": eventId,
			},
		})
	}

	return resp, nil
}
