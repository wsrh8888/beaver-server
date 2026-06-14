package logic

import (
	"context"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsBlockedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsBlockedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsBlockedLogic {
	return &IsBlockedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsBlockedLogic) IsBlocked(in *friend_rpc.IsBlockedReq) (*friend_rpc.IsBlockedRes, error) {
	var blockCount int64
	l.svcCtx.DB.Model(&friend_models.FriendBlockModel{}).
		Where("(user_id = ? AND blocked_user_id = ?) OR (user_id = ? AND blocked_user_id = ?)",
			in.UserA, in.UserB, in.UserB, in.UserA).
		Count(&blockCount)
	return &friend_rpc.IsBlockedRes{IsBlocked: blockCount > 0}, nil
}
