package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteFriendsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除好友关系
func NewBatchDeleteFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteFriendsLogic {
	return &BatchDeleteFriendsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteFriendsLogic) BatchDeleteFriends(req *types.BatchDeleteFriendsReq) (resp *types.BatchDeleteFriendsRes, err error) {
	// 转换字符串UUID为string切片（现在使用UUID而不是数据库ID）
	var friendUUIDs []string
	for _, uuidStr := range req.Ids {
		// 简单验证UUID格式（可以根据需要添加更严格的验证）
		if len(uuidStr) == 0 {
			logx.Errorf("无效的好友关系UUID: %s", uuidStr)
			return nil, fmt.Errorf("无效的好友关系UUID: %s", uuidStr)
		}
		friendUUIDs = append(friendUUIDs, uuidStr)
	}

	// 先查询要删除的好友关系
	var friends []friend_models.FriendModel
	err = l.svcCtx.DB.Where("uuid IN ?", friendUUIDs).Find(&friends).Error
	if err != nil {
		logx.Errorf("查询要删除的好友关系失败: %v", err)
		return nil, err
	}

	if len(friends) == 0 {
		return nil, errors.New("没有找到要删除的好友关系")
	}

	// 批量删除好友关系（物理删除）
	err = l.svcCtx.DB.Unscoped().Where("uuid IN ?", friendUUIDs).Delete(&friend_models.FriendModel{}).Error
	if err != nil {
		logx.Errorf("批量删除好友关系失败: %v", err)
		return nil, err
	}

	logx.Infof("批量删除好友关系完成, 删除数量: %d", len(friends))
	for _, friend := range friends {
		logx.Infof("删除好友关系 - UUID: %s, SendUserID: %s, RevUserID: %s",
			friend.UUID, friend.SendUserID, friend.RevUserID)
	}

	return &types.BatchDeleteFriendsRes{}, nil
}
