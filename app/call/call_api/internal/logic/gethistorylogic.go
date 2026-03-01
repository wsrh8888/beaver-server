package logic

import (
	"context"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询我的通话历史记录
func NewGetHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHistoryLogic {
	return &GetHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHistoryLogic) GetHistory(req *types.CallHistoryReq) (resp *types.CallHistoryRes, err error) {
	var sessions []call_models.CallSession

	// 简单实现：查询参与者包含该用户的会话
	// 实际开发中通常通过联合查询或参与者表反查
	query := l.svcCtx.DB.Table("call_sessions").
		Joins("JOIN call_participants ON call_sessions.room_id = call_participants.room_id").
		Where("call_participants.user_id = ?", req.UserID).
		Order("call_sessions.created_at DESC")

	var total int64
	query.Count(&total)

	err = query.Offset((req.Page - 1) * req.Size).Limit(req.Size).Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	list := make([]types.CallHistoryResItem, 0, len(sessions))
	for _, s := range sessions {
		startTime := int64(0)
		if s.StartTime != nil {
			startTime = s.StartTime.Unix()
		}
		list = append(list, types.CallHistoryResItem{
			RoomID:    s.RoomID,
			CallerID:  s.CallerID,
			CallType:  s.CallType,
			Status:    int8(s.Status),
			StartTime: startTime,
			Duration:  s.Duration,
		})
	}

	return &types.CallHistoryRes{
		List: list,
	}, nil
}
