package logic

import (
	"context"
	"time"

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadByEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按事件ID标记单个通知已读
func NewMarkReadByEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadByEventLogic {
	return &MarkReadByEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkReadByEventLogic) MarkReadByEvent(req *types.MarkReadByEventReq) (resp *types.MarkReadByEventRes, err error) {
	userId := req.UserID
	eventId := req.EventID

	// 更新指定通知为已读
	result := l.svcCtx.DB.Model(&notification_models.NotificationInbox{}).
		Where("user_id = ? AND event_id = ? AND is_deleted = ?", userId, eventId, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": time.Now(),
		})

	if result.Error != nil {
		l.Logger.Errorf("标记单个通知已读失败: %v", result.Error)
		return nil, result.Error
	}

	resp = &types.MarkReadByEventRes{}

	return resp, nil
}
