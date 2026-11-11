package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ListNotificationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotificationsLogic {
	return &ListNotificationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListNotificationsLogic", ctx),
	}
}

func (l *ListNotificationsLogic) ListNotifications(in *doc_rpc.ListNotificationsReq) (*doc_rpc.ListNotificationsRes, error) {
	uid := trim(in.UserId)
	if uid == "" {
		return &doc_rpc.ListNotificationsRes{}, nil
	}
	q := l.svcCtx.DB.Model(&doc_models.DocNotification{}).Where("user_id = ?", uid)
	var unread int64
	if err := l.svcCtx.DB.Model(&doc_models.DocNotification{}).
		Where("user_id = ? AND is_read = ?", uid, false).Count(&unread).Error; err != nil {
		return nil, err
	}
	if in.UnreadOnly {
		q = q.Where("is_read = ?", false)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}
	page, size := pageWindow(in.Page, in.PageSize)
	var rows []doc_models.DocNotification
	if err := q.Order("created_at_str desc, id desc").
		Offset((page - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]*doc_rpc.NotificationItem, 0, len(rows))
	for i := range rows {
		items = append(items, toNotificationItem(&rows[i]))
	}
	return &doc_rpc.ListNotificationsRes{Total: total, Unread: unread, Items: items}, nil
}
