package logic

import (
	"context"

	"beaver/app/notification/notification_models"
	"beaver/app/notification/notification_rpc/internal/svc"
	"beaver/app/notification/notification_rpc/types/notification_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetInboxVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInboxVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInboxVersionsLogic {
	return &GetInboxVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetInboxVersionsLogic) GetInboxVersions(in *notification_rpc.GetInboxVersionsReq) (*notification_rpc.GetInboxVersionsRes, error) {
	resp := &notification_rpc.GetInboxVersionsRes{
		InboxVersions: []*notification_rpc.InboxVersion{},
		MaxVersion:    0,
	}

	if in.UserId == "" {
		return resp, nil
	}

	var rows []notification_models.NotificationInbox
	query := l.svcCtx.DB.WithContext(l.ctx).
		Where("user_id = ? AND version > ?", in.UserId, in.SinceVersion).
		Order("version ASC")
	if in.Limit > 0 {
		query = query.Limit(int(in.Limit))
	}

	if err := query.Find(&rows).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	for _, row := range rows {
		resp.InboxVersions = append(resp.InboxVersions, &notification_rpc.InboxVersion{
			EventId: row.EventID,
			Version: row.Version,
		})
		if row.Version > resp.MaxVersion {
			resp.MaxVersion = row.Version
		}
	}

	return resp, nil
}
