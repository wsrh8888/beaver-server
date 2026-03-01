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

// 核心：创建通话会话并初始化参与者名单
func (l *CreateSessionLogic) CreateSession(in *call_rpc.CreateSessionReq) (*call_rpc.CreateSessionRes, error) {
	var participants []*call_rpc.Participant

	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建会话
		session := &call_models.CallSession{
			RoomID:         in.RoomId,
			CallerID:       in.CallerId,
			CallType:       int8(in.CallType),
			ConversationID: in.ConversationId,                // 存入会话ID
			Status:         call_models.SessionStatusCalling, // 1-进行中
		}
		if err := tx.Create(session).Error; err != nil {
			return err
		}

		// 2. 创建参与者 (发起者)
		caller := &call_models.CallParticipant{
			RoomID: in.RoomId,
			UserID: in.CallerId,
			Status: call_models.ParticipantStatusJoined, // 2-已接听 (发起者默认已进入)
			Role:   1,                                   // 1-发起者
		}
		if err := tx.Create(caller).Error; err != nil {
			return err
		}
		participants = append(participants, &call_rpc.Participant{
			UserId: in.CallerId,
			Status: int32(call_models.ParticipantStatusJoined),
		})

		// 3. 创建参与者 (受邀者) - 仅单聊初始化对方，群聊在邀请或主动加入时处理
		if in.CallType == 1 {
			target := &call_models.CallParticipant{
				RoomID: in.RoomId,
				UserID: in.TargetId,
				Status: call_models.ParticipantStatusCalling, // 1-待接听
				Role:   2,                                    // 2-受邀者
			}
			if err := tx.Create(target).Error; err != nil {
				return err
			}
			participants = append(participants, &call_rpc.Participant{
				UserId: in.TargetId,
				Status: int32(call_models.ParticipantStatusCalling),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &call_rpc.CreateSessionRes{
		Success:      true,
		Participants: participants,
	}, nil
}
