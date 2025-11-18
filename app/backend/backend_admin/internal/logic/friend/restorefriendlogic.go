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

type RestoreFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 恢复好友关系
func NewRestoreFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreFriendLogic {
	return &RestoreFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RestoreFriendLogic) RestoreFriend(req *types.RestoreFriendReq) (resp *types.RestoreFriendRes, err error) {
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

	// 检查是否已经是未删除状态
	if !friend.IsDeleted {
		logx.Infof("好友关系已经是未删除状态, Id: %s", req.FriendID)
		return &types.RestoreFriendRes{}, nil
	}

	// 恢复好友关系（设置IsDeleted为false）
	err = l.svcCtx.DB.Model(&friend).Update("is_deleted", false).Error
	if err != nil {
		logx.Errorf("恢复好友关系失败: %v", err)
		return nil, err
	}

	logx.Infof("好友关系恢复成功, Id: %s, SendUserID: %s, RevUserID: %s",
		req.FriendID, friend.SendUserID, friend.RevUserID)
	return &types.RestoreFriendRes{}, nil
}
