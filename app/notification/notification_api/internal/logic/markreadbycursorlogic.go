package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type MarkReadByCursorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按分类游标标记已读（高效批量，首版主路径）
func NewMarkReadByCursorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadByCursorLogic {
	return &MarkReadByCursorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkReadByCursorLogic) MarkReadByCursor(req *types.MarkReadByCursorReq) (resp *types.MarkReadByCursorRes, err error) {
	resp = &types.MarkReadByCursorRes{Affected: 0}

	if req.UserID == "" || req.Category == "" {
		return resp, errors.New("userId 和 category 不能为空")
	}
	if req.ToVersion <= 0 && req.ToEventID == "" {
		return resp, errors.New("需提供 toVersion 或 toEventId")
	}

	now := time.Now()

	err = l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 生成按用户+分类递增的光标版本（通过组合键）
		cursorKey := req.UserID + "_" + req.Category
		cursorVersion := l.svcCtx.VersionGen.GetNextVersion(notification_models.VersionScopeCursorPerUser, "user_category", cursorKey)
		if cursorVersion == -1 {
			return errors.New("生成游标版本失败")
		}

		query := tx.Model(&notification_models.NotificationInbox{}).
			Where("user_id = ? AND category = ?", req.UserID, req.Category)

		if req.ToVersion > 0 {
			query = query.Where("version <= ?", req.ToVersion)
		} else {
			query = query.Where("event_id = ?", req.ToEventID)
		}

		update := map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}

		result := query.Updates(update)
		if result.Error != nil {
			return result.Error
		}
		resp.Affected = result.RowsAffected

		// upsert read cursor
		var cursor notification_models.NotificationReadCursor
		err := tx.Where("user_id = ? AND category = ?", req.UserID, req.Category).
			First(&cursor).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cursor = notification_models.NotificationReadCursor{
				UserID:       req.UserID,
				Category:     req.Category,
				Version:      cursorVersion,
				LastEventID:  req.ToEventID,
				LastReadAt:   &now,
				LastReadTime: now.UnixMilli(),
			}
			return tx.Create(&cursor).Error
		} else if err != nil {
			return err
		}

		updateCursor := map[string]interface{}{
			"version":        cursorVersion,
			"last_event_id":  req.ToEventID,
			"last_read_at":   now,
			"last_read_time": now.UnixMilli(),
		}
		return tx.Model(&cursor).Updates(updateCursor).Error
	})

	return resp, err
}
