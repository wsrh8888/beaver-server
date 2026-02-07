package logic

import (
	"context"
	"time"

	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/internal/svc"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FinalizeSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFinalizeSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FinalizeSessionLogic {
	return &FinalizeSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 结束通话记录
func (l *FinalizeSessionLogic) FinalizeSession(in *call_rpc.FinalizeSessionReq) (*call_rpc.FinalizeSessionRes, error) {
	now := time.Now()
	err := l.svcCtx.DB.Model(&call_models.CallSession{}).
		Where("room_id = ?", in.RoomId).
		Updates(map[string]interface{}{
			"status":   int8(in.Status),
			"end_time": &now,
			"duration": in.Duration,
		}).Error

	if err != nil {
		return nil, err
	}

	return &call_rpc.FinalizeSessionRes{Success: true}, nil
}
