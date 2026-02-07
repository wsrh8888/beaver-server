package logic

import (
	"context"

	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/internal/svc"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionLogic {
	return &GetSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取通话信息
func (l *GetSessionLogic) GetSession(in *call_rpc.GetSessionReq) (*call_rpc.GetSessionRes, error) {
	var session call_models.CallSession
	if err := l.svcCtx.DB.Where("room_id = ?", in.RoomId).First(&session).Error; err != nil {
		return nil, err
	}

	var participantIds []string
	l.svcCtx.DB.Model(&call_models.CallParticipant{}).
		Where("room_id = ?", in.RoomId).
		Pluck("user_id", &participantIds)

	return &call_rpc.GetSessionRes{
		RoomId:         session.RoomID,
		CallerId:       session.CallerID,
		CallType:       int32(session.CallType),
		Status:         int32(session.Status),
		ParticipantIds: participantIds,
	}, nil
}
