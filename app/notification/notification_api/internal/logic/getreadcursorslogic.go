package logic

import (
	"context"

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetReadCursorsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按分类拉取通知已读游标
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

	var rows []notification_models.NotificationRead
	query := l.svcCtx.DB.WithContext(l.ctx).
		Where("user_id = ?", req.UserID)

	if len(req.Categories) > 0 {
		query = query.Where("category IN ?", req.Categories)
	}

	if err = query.Find(&rows).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	for _, row := range rows {
		lastReadAt := int64(0)
		if row.LastReadAt != nil {
			lastReadAt = row.LastReadAt.UnixMilli()
		}

		resp.Cursors = append(resp.Cursors, types.ReadCursorItem{
			Category:   row.Category,
			Version:    row.Version,
			LastReadAt: lastReadAt,
		})
	}

	return resp, nil
}
