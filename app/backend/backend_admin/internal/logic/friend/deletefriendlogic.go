package logic

import (
	"context"
	"errors"
	"strconv"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 强制删除好友关系
func NewDeleteFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendLogic {
	return &DeleteFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFriendLogic) DeleteFriend(req *types.DeleteFriendReq) (resp *types.DeleteFriendRes, err error) {
	// 转换ID
	friendID, err := strconv.ParseUint(req.FriendID, 10, 32)
	if err != nil {
		logx.Errorf("无效的好友关系ID: %s", req.FriendID)
		return nil, errors.New("无效的好友关系ID")
	}

	// 先查询好友关系是否存在
	var friend friend_models.FriendModel
	err = l.svcCtx.DB.Where("id = ?", uint(friendID)).First(&friend).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("好友关系不存在, Id: %s", req.FriendID)
			return nil, errors.New("好友关系不存在")
		}
		logx.Errorf("查询好友关系失败: %v", err)
		return nil, err
	}

	// 强制删除好友关系（物理删除）
	err = l.svcCtx.DB.Unscoped().Delete(&friend).Error
	if err != nil {
		logx.Errorf("删除好友关系失败: %v", err)
		return nil, err
	}

	logx.Infof("好友关系删除成功, Id: %s, SendUserID: %s, RevUserID: %s",
		req.FriendID, friend.SendUserID, friend.RevUserID)
	return &types.DeleteFriendRes{}, nil
}
