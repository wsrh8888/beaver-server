package logic

import (
	"context"

	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/internal/svc"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetParticipantsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetParticipantsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetParticipantsLogic {
	return &GetParticipantsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取参与者列表及状态
func (l *GetParticipantsLogic) GetParticipants(in *call_rpc.GetParticipantsReq) (*call_rpc.GetParticipantsRes, error) {
	var participants []call_models.CallParticipant
	if err := l.svcCtx.DB.Where("room_id = ?", in.RoomId).Find(&participants).Error; err != nil {
		return nil, err
	}

	var res []*call_rpc.Participant
	for _, p := range participants {
		res = append(res, &call_rpc.Participant{
			UserId: p.UserID,
			Status: int32(p.Status),
		})
	}

	return &call_rpc.GetParticipantsRes{
		Participants: res,
	}, nil
}
