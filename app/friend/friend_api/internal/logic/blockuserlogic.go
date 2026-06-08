package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)


type BlockUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// 拉黑用户
func NewBlockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockUserLogic {
	return &BlockUserLogic{
		ctx:    ctx,
		logger: logger.New("block_user"),
		svcCtx: svcCtx,
	}
}

func (l *BlockUserLogic) BlockUser(req *types.BlockUserReq) (resp *types.BlockUserRes, err error) {
	if req.UserID == req.BlockedUserID {
		return nil, errors.New("不能拉黑自己")
	}

	// 检查是否已拉黑
	var existing friend_models.FriendBlockModel
	result := l.svcCtx.DB.Where("user_id = ? AND blocked_user_id = ?", req.UserID, req.BlockedUserID).First(&existing)
	if result.Error == nil {
		return &types.BlockUserRes{}, nil // 已拉黑，幂等处理
	}

	block := friend_models.FriendBlockModel{
		BlockID:       uuid.New().String(),
		UserID:        req.UserID,
		BlockedUserID: req.BlockedUserID,
	}
	if err = l.svcCtx.DB.Create(&block).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("拉黑用户失败: userID=%s blockedUserID=%s err=%v", req.UserID, req.BlockedUserID, err)
		return nil, errors.New("操作失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "拉黑用户成功",
		Data: map[string]interface{}{
			"userId":        req.UserID,
			"blockedUserId": req.BlockedUserID,
		},
	})

	return &types.BlockUserRes{}, nil
}
