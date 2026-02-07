package logic

import (
	"context"

	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/internal/svc"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 供 Call-Api 调用，创建通话记录
func (l *CreateSessionLogic) CreateSession(in *call_rpc.CreateSessionReq) (*call_rpc.CreateSessionRes, error) {
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 创建会话
		session := &call_models.CallSession{
			RoomID:   in.RoomId,
			CallerID: in.CallerId,
			CallType: int8(in.CallType),
			Status:   1, // 1-呼叫中
		}
		if err := tx.Create(session).Error; err != nil {
			return err
		}

		// 创建参与者 (发起者)
		caller := &call_models.CallParticipant{
			RoomID: in.RoomId,
			UserID: in.CallerId,
			Status: 2, // 2-已接听 (发起者默认已进入)
			Role:   1, // 1-发起者
		}
		if err := tx.Create(caller).Error; err != nil {
			return err
		}

		// 创建参与者 (受邀者)
		target := &call_models.CallParticipant{
			RoomID: in.RoomId,
			UserID: in.TargetId,
			Status: 1, // 1-待接听
			Role:   2, // 2-受邀者
		}
		if err := tx.Create(target).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &call_rpc.CreateSessionRes{Success: true}, nil
}
