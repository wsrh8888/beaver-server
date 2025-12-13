package logic

import (
	"context"
	"time"

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadByCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按分类标记所有通知为已读
func NewMarkReadByCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadByCategoryLogic {
	return &MarkReadByCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkReadByCategoryLogic) MarkReadByCategory(req *types.MarkReadByCategoryReq) (resp *types.MarkReadByCategoryRes, err error) {
	userId := req.UserID
	category := req.Category

	now := time.Now()

	// 更新该用户该分类下所有未读通知为已读
	result := l.svcCtx.DB.Model(&notification_models.NotificationInbox{}).
		Where("user_id = ? AND category = ? AND is_read = ? AND is_deleted = ?",
			userId, category, false, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		})

	if result.Error != nil {
		l.Logger.Errorf("标记分类已读失败: %v", result.Error)
		return nil, result.Error
	}

	// 更新游标记录
	cursor := &notification_models.NotificationRead{
		UserID:     userId,
		Category:   category,
		Version:    1, // 简化版本管理
		LastReadAt: &now,
	}

	err = l.svcCtx.DB.Where("user_id = ? AND category = ?", userId, category).
		Assign(cursor).
		FirstOrCreate(cursor).Error

	if err != nil {
		l.Logger.Errorf("更新游标失败: %v", err)
		return nil, err
	}

	resp = &types.MarkReadByCategoryRes{
		Affected: result.RowsAffected,
	}

	return resp, nil
}
