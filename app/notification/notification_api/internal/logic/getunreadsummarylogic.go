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

	// 如果没有指定分类，默认查询所有分类
	if len(cats) == 0 {
		cats = []string{notification_models.CategorySocial, notification_models.CategoryGroup, notification_models.CategoryMoment}
	}

	// 为每个分类计算结合游标时间的未读数
	for _, category := range cats {
		// 获取该分类的游标时间
		var cursor notification_models.NotificationRead
		err := l.svcCtx.DB.WithContext(l.ctx).
			Where("user_id = ? AND category = ?", req.UserID, category).
			First(&cursor).Error

		var lastReadAt int64
		if err == nil && cursor.LastReadAt != nil {
			lastReadAt = cursor.LastReadAt.Unix()
		}

		// 计算未读数：created_at > last_read_at 且 is_read = false
		var unreadCount int64
		query := l.svcCtx.DB.WithContext(l.ctx).
			Model(&notification_models.NotificationInbox{}).
			Where("user_id = ? AND category = ? AND is_read = ? AND is_deleted = ?",
				req.UserID, category, false, false)

		// 如果有游标时间，只统计游标时间之后的未读通知
		if lastReadAt > 0 {
			query = query.Where("created_at > ?", lastReadAt)
		}

		if err = query.Count(&unreadCount).Error; err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		if unreadCount > 0 {
			resp.ByCat = append(resp.ByCat, types.CategoryUnreadItem{
				Category: category,
				Unread:   unreadCount,
			})
			resp.Total += unreadCount
		}
	}

	return resp, nil
}
