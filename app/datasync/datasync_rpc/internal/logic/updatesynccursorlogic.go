package logic

import (
	"context"

	sync_models "beaver/app/datasync/datasync_models"
	"beaver/app/datasync/datasync_rpc/internal/svc"
	"beaver/app/datasync/datasync_rpc/types/types/datasync_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSyncCursorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSyncCursorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSyncCursorLogic {
	return &UpdateSyncCursorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新同步游标
func (l *UpdateSyncCursorLogic) UpdateSyncCursor(in *datasync_rpc.UpdateSyncCursorReq) (*datasync_rpc.UpdateSyncCursorRes, error) {
	var cursor sync_models.DatasyncModel

	// 先查询是否存在记录
	err := l.svcCtx.DB.Where("user_id = ? AND device_id = ? AND data_type = ?",
		in.UserId, in.DeviceId, in.DataType).First(&cursor).Error

	if err != nil {
		// 如果没有找到记录，创建新记录
		if err.Error() == "record not found" {
			cursor = sync_models.DatasyncModel{
				UserID:         in.UserId,
				DeviceID:       in.DeviceId,
				DataType:       in.DataType,
				ConversationID: in.ConversationId,
				LastSeq:        in.LastSeq, // 使用LastSeq替代LastSyncTime
				SyncStatus:     "completed",
			}
			err = l.svcCtx.DB.Create(&cursor).Error
		} else {
			l.Errorf("查询同步游标失败: %v", err)
			return nil, err
		}
	} else {
		// 更新现有记录
		cursor.ConversationID = in.ConversationId
		cursor.LastSeq = in.LastSeq // 使用LastSeq替代LastSyncTime
		cursor.SyncStatus = "completed"
		err = l.svcCtx.DB.Save(&cursor).Error
	}

	if err != nil {
		l.Errorf("更新同步游标失败: %v", err)
		return nil, err
	}

	return &datasync_rpc.UpdateSyncCursorRes{}, nil
}
