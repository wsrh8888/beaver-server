package logic

import (
	"context"
	"errors"

	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/internal/svc"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserStatusLogic {
	return &GetUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 供 Chat/User 等服务查询用户是否忙碌
func (l *GetUserStatusLogic) GetUserStatus(in *call_rpc.GetUserStatusReq) (*call_rpc.GetUserStatusRes, error) {
	var participant call_models.CallParticipant
	err := l.svcCtx.DB.Where("user_id = ? AND status IN ?", in.UserId, []int8{1, 2}).First(&participant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &call_rpc.GetUserStatusRes{
				IsBusy: false,
			}, nil
		}
		return nil, err
	}

	return &call_rpc.GetUserStatusRes{
		IsBusy: true,
		RoomId: participant.RoomID,
	}, nil
}
