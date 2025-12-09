package logic

import (
	"context"

	"beaver/app/notification/notification_models"
	"beaver/app/notification/notification_rpc/internal/svc"
	"beaver/app/notification/notification_rpc/types/notification_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetReadCursorVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetReadCursorVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReadCursorVersionsLogic {
	return &GetReadCursorVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetReadCursorVersionsLogic) GetReadCursorVersions(in *notification_rpc.GetReadCursorVersionsReq) (*notification_rpc.GetReadCursorVersionsRes, error) {
	resp := &notification_rpc.GetReadCursorVersionsRes{
		CursorVersions: []*notification_rpc.ReadCursorVersion{},
		MaxVersion:     0,
	}

	if in.UserId == "" {
		return resp, nil
	}

	var rows []notification_models.NotificationReadCursor
	query := l.svcCtx.DB.WithContext(l.ctx).
		Where("user_id = ? AND version > ?", in.UserId, in.SinceVersion).
		Order("version ASC")

	if err := query.Find(&rows).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	for _, row := range rows {
		resp.CursorVersions = append(resp.CursorVersions, &notification_rpc.ReadCursorVersion{
			Category: row.Category,
			Version:  row.Version,
		})
		if row.Version > resp.MaxVersion {
			resp.MaxVersion = row.Version
		}
	}

	return resp, nil
}
