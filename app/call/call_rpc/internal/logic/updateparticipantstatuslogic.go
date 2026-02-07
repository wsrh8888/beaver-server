package logic

import (
	"context"
	"time"

	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/internal/svc"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateParticipantStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateParticipantStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateParticipantStatusLogic {
	return &UpdateParticipantStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新参与者状态
func (l *UpdateParticipantStatusLogic) UpdateParticipantStatus(in *call_rpc.UpdateParticipantStatusReq) (*call_rpc.UpdateParticipantStatusRes, error) {
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		var updates map[string]interface{}
		if in.Status == 2 { // 2-已接听
			now := time.Now()
			updates = map[string]interface{}{
				"status":    int8(in.Status),
				"join_time": &now,
			}
		} else {
			updates = map[string]interface{}{
				"status": int8(in.Status),
			}
		}

		err := tx.Model(&call_models.CallParticipant{}).
			Where("room_id = ? AND user_id = ?", in.RoomId, in.UserId).
			Updates(updates).Error
		if err != nil {
			return err
		}

		// 如果状态是已接听，也将 session 状态改为进行中
		if in.Status == 2 {
			now := time.Now()
			err = tx.Model(&call_models.CallSession{}).
				Where("room_id = ?", in.RoomId).
				Updates(map[string]interface{}{
					"status":     2, // 2-进行中
					"start_time": &now,
				}).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &call_rpc.UpdateParticipantStatusRes{Success: true}, nil
}
