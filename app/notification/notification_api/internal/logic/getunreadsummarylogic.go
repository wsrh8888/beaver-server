package logic

import (
	"context"

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUnreadSummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取未读汇总（红点）
func NewGetUnreadSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnreadSummaryLogic {
	return &GetUnreadSummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUnreadSummaryLogic) GetUnreadSummary(req *types.GetUnreadSummaryReq) (resp *types.GetUnreadSummaryRes, err error) {
	resp = &types.GetUnreadSummaryRes{Total: 0, ByCat: []types.CategoryUnreadItem{}}

	if req.UserID == "" {
		return resp, nil
	}

	// 使用 scenes 作为类别过滤兜底
	cats := req.Categories
	if len(cats) == 0 && len(req.Scenes) > 0 {
		cats = req.Scenes
	}

	type aggRow struct {
		Category string
		Unread   int64
	}
	var rows []aggRow

	query := l.svcCtx.DB.WithContext(l.ctx).
		Model(&notification_models.NotificationInbox{}).
		Select("category, COUNT(*) as unread").
		Where("user_id = ? AND is_read = ?", req.UserID, false).
		Group("category")

	if len(cats) > 0 {
		query = query.Having("category IN ?", cats)
	}

	if err = query.Find(&rows).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	for _, row := range rows {
		resp.ByCat = append(resp.ByCat, types.CategoryUnreadItem{
			Category: row.Category,
			Unread:   row.Unread,
		})
		resp.Total += row.Unread
	}

	return resp, nil
}
