package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNoticeUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeUpdateLogic {
	return &NoticeUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeUpdateLogic) NoticeUpdate(req *types.NoticeUpdateReq) (resp *types.NoticeUpdateRes, err error) {
	// 参数验证
	if req.UserID == "" || req.FriendID == "" {
		return nil, errors.New("用户ID和好友ID不能为空")
	}

	// 不能修改自己的备注
	if req.UserID == req.FriendID {
		return nil, errors.New("不能修改自己的备注")
	}

	var friend friend_models.FriendModel

	// 检查是否为好友关系
	if !friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID) {
		l.Logger.Errorf("尝试修改非好友备注: userID=%s, friendID=%s", req.UserID, req.FriendID)
		return nil, errors.New("不是好友关系")
	}

	// 查询好友关系详情
	err = l.svcCtx.DB.Take(&friend,
		"((send_user_id = ? AND rev_user_id = ?) OR (send_user_id = ? AND rev_user_id = ?)) AND is_deleted = 0",
		req.UserID, req.FriendID, req.FriendID, req.UserID).Error
	if err != nil {
		l.Logger.Errorf("查询好友关系失败: %v", err)
		return nil, errors.New("查询好友关系失败")
	}

	// 根据用户角色更新对应的备注字段
	if friend.SendUserID == req.UserID {
		// 我是发起方，更新发起方备注
		if friend.SendUserNotice == req.Notice {
			// 备注没有变化，直接返回
			return &types.NoticeUpdateRes{}, nil
		}
		err = l.svcCtx.DB.Model(&friend_models.FriendModel{}).Where("id = ?", friend.ID).Update("send_user_notice", req.Notice).Error
	} else if friend.RevUserID == req.UserID {
		// 我是接收方，更新接收方备注
		if friend.RevUserNotice == req.Notice {
			// 备注没有变化，直接返回
			return &types.NoticeUpdateRes{}, nil
		}
		err = l.svcCtx.DB.Model(&friend_models.FriendModel{}).Where("id = ?", friend.ID).Update("rev_user_notice", req.Notice).Error
	} else {
		// 这种情况理论上不应该发生
		l.Logger.Errorf("用户角色异常: userID=%s, friendID=%s", req.UserID, req.FriendID)
		return nil, errors.New("用户角色异常")
	}

	if err != nil {
		l.Logger.Errorf("更新好友备注失败: %v", err)
		return nil, errors.New("更新好友备注失败")
	}

	l.Logger.Infof("更新好友备注成功: userID=%s, friendID=%s, notice=%s", req.UserID, req.FriendID, req.Notice)
	return &types.NoticeUpdateRes{}, nil
}
