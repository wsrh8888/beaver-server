package logic

import (
	"context"
	"errors"
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

	// 生成游标版本号（按用户+分类分区）
	cursorVersion := l.svcCtx.VersionGen.GetNextVersion(notification_models.VersionScopeCursorPerUser, "user_id", userId)
	if cursorVersion == -1 {
		l.Logger.Errorf("生成游标版本号失败")
		return nil, errors.New("生成版本号失败")
	}

	// 更新或创建游标记录：优先更新，不存在则创建
	result := l.svcCtx.DB.Model(&notification_models.NotificationRead{}).
		Where("user_id = ? AND category = ?", userId, category).
		Updates(map[string]interface{}{
			"version":      cursorVersion,
			"last_read_at": now,
			"updated_at":   now,
		})

	if result.Error != nil {
		l.Logger.Errorf("更新游标失败: %v", result.Error)
		return nil, result.Error
	}

	// 如果没有更新到记录，说明记录不存在，需要创建
	if result.RowsAffected == 0 {
		cursor := &notification_models.NotificationRead{
			UserID:     userId,
			Category:   category,
			Version:    cursorVersion,
			LastReadAt: &now,
		}

		err = l.svcCtx.DB.Create(cursor).Error
		if err != nil {
			l.Logger.Errorf("创建游标失败: %v", err)
			return nil, err
		}
	}

	resp = &types.MarkReadByCategoryRes{}

	return resp, nil
}
