package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type UnblockUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// 取消拉黑
func NewUnblockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnblockUserLogic {
	return &UnblockUserLogic{
		ctx:    ctx,
		logger: logger.New("unblock_user"),
		svcCtx: svcCtx,
	}
}

func (l *UnblockUserLogic) UnblockUser(req *types.UnblockUserReq) (resp *types.UnblockUserRes, err error) {
	result := l.svcCtx.DB.Where("user_id = ? AND blocked_user_id = ?", req.UserID, req.BlockedUserID).
		Delete(&friend_models.FriendBlockModel{})
	if result.Error != nil {
		logx.WithContext(l.ctx).Errorf("取消拉黑失败: userID=%s blockedUserID=%s err=%v", req.UserID, req.BlockedUserID, result.Error)
		return nil, errors.New("操作失败")
	}
	l.logger.Info(model.LogMsg{
		Text: "取消拉黑成功",
		Data: map[string]interface{}{
			"userId":        req.UserID,
			"blockedUserId": req.BlockedUserID,
		},
	})
	return &types.UnblockUserRes{}, nil
}
