package logic

import (
	"context"
	"time"

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReadCursorsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按分类拉取已读游标
func NewGetReadCursorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReadCursorsLogic {
	return &GetReadCursorsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReadCursorsLogic) GetReadCursors(req *types.GetReadCursorsReq) (resp *types.GetReadCursorsRes, err error) {
	resp = &types.GetReadCursorsRes{Cursors: []types.ReadCursorItem{}}

	if req.UserID == "" {
		return resp, nil
	}

	query := l.svcCtx.DB.WithContext(l.ctx).Model(&notification_models.NotificationReadCursor{}).
		Where("user_id = ?", req.UserID)
	if len(req.Categories) > 0 {
		query = query.Where("category IN ?", req.Categories)
	}

	var rows []notification_models.NotificationReadCursor
	if err = query.Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		resp.Cursors = append(resp.Cursors, types.ReadCursorItem{
			Category:     row.Category,
			Version:      row.Version,
			LastEventID:  row.LastEventID,
			LastReadTime: row.LastReadTime,
			LastReadAt:   time.Time(*row.LastReadAt).UnixMilli(),
		})
	}

	return resp, nil
}
