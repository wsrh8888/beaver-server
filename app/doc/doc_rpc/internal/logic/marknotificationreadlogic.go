package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type MarkNotificationReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewMarkNotificationReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkNotificationReadLogic {
	return &MarkNotificationReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("MarkNotificationReadLogic", ctx),
	}
}

func (l *MarkNotificationReadLogic) MarkNotificationRead(in *doc_rpc.MarkNotificationReadReq) (*doc_rpc.MarkNotificationReadRes, error) {
	uid := trim(in.UserId)
	if uid == "" {
		return nil, errInvalidArg
	}
	db := l.svcCtx.DB.Model(&doc_models.DocNotification{}).Where("user_id = ?", uid)
	if in.All {
		if err := db.Where("is_read = ?", false).Update("is_read", true).Error; err != nil {
			return nil, err
		}
	} else if nid := trim(in.NotificationId); nid != "" {
		if err := db.Where("notification_id = ?", nid).Update("is_read", true).Error; err != nil {
			return nil, err
		}
	} else {
		return nil, errInvalidArg
	}
	var unread int64
	if err := l.svcCtx.DB.Model(&doc_models.DocNotification{}).
		Where("user_id = ? AND is_read = ?", uid, false).Count(&unread).Error; err != nil {
		return nil, err
	}
	return &doc_rpc.MarkNotificationReadRes{Unread: unread}, nil
}
