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
		now := time.Now()

		// 1. 查找现有记录
		var p call_models.CallParticipant
		err := tx.Where("room_id = ? AND user_id = ?", in.RoomId, in.UserId).First(&p).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		// 2. 准备更新的数据
		updates := make(map[string]interface{})
		updates["status"] = in.Status

		if in.Status == int32(call_models.ParticipantStatusJoined) { // 2-已接听 (joined)
			updates["join_time"] = &now
		} else if in.Status >= int32(call_models.ParticipantStatusRejected) && in.Status <= int32(call_models.ParticipantStatusLeft) { // 3-拒绝, 4-超时, 5-挂断
			updates["leave_time"] = &now
		}

		if err == gorm.ErrRecordNotFound {
			// 如果不存在且是由于邀请或加入产生，则新建
			var session call_models.CallSession
			tx.Where("room_id = ?", in.RoomId).First(&session)

			role := int8(2) // 默认受邀者
			if session.CallerID == in.UserId {
				role = 1 // 发起者
			}

			p = call_models.CallParticipant{
				RoomID: in.RoomId,
				UserID: in.UserId,
				Status: call_models.ParticipantStatus(in.Status),
				Role:   role,
			}
			if in.Status == int32(call_models.ParticipantStatusJoined) {
				p.JoinTime = &now
			}
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
		} else {
			// 存在则更新
			if err := tx.Model(&p).Updates(updates).Error; err != nil {
				return err
			}
		}

		// 3. 如果状态是"已接听" (2)，也将 session 状态改为进行中 (1)
		if in.Status == int32(call_models.ParticipantStatusJoined) {
			err = tx.Model(&call_models.CallSession{}).
				Where("room_id = ?", in.RoomId).
				Updates(map[string]interface{}{
					"status":     call_models.SessionStatusCalling, // 1-进行中
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
