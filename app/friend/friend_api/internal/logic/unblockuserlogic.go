package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnblockUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消拉黑
func NewUnblockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnblockUserLogic {
	return &UnblockUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnblockUserLogic) UnblockUser(req *types.UnblockUserReq) (resp *types.UnblockUserRes, err error) {
	result := l.svcCtx.DB.Where("user_id = ? AND blocked_user_id = ?", req.UserID, req.BlockedUserID).
		Delete(&friend_models.FriendBlockModel{})
	if result.Error != nil {
		l.Errorf("取消拉黑失败: userID=%s blockedUserID=%s err=%v", req.UserID, req.BlockedUserID, result.Error)
		return nil, errors.New("操作失败")
	}
	return &types.UnblockUserRes{}, nil
}
