package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendLogic {
	return &AddFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFriendLogic) AddFriend(req *types.AddFriendReq) (resp *types.AddFriendRes, err error) {
	var friend friend_models.FriendModel

	// 检查是否已经是好友
	if friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID) {
		return nil, errors.New("已经是好友了")
	}

	// 检查目标用户是否存在
	var userInfo user_models.UserModel
	err = l.svcCtx.DB.Take(&userInfo, "uuid = ?", req.FriendID).Error
	if err != nil {
		l.Logger.Errorf("目标用户不存在: friendID=%s, error=%v", req.FriendID, err)
		return nil, errors.New("用户不存在")
	}

	// 检查是否已经有待处理的好友请求
	var existingVerify friend_models.FriendVerifyModel
	err = l.svcCtx.DB.Take(&existingVerify,
		"(send_user_id = ? AND rev_user_id = ? AND rev_status = 0) OR (send_user_id = ? AND rev_user_id = ? AND rev_status = 0)",
		req.UserID, req.FriendID, req.FriendID, req.UserID).Error

	if err == nil {
		l.Logger.Infof("已存在待处理的好友请求: userID=%s, friendID=%s", req.UserID, req.FriendID)
		return &types.AddFriendRes{}, nil
	}

	// 创建好友验证请求
	verifyModel := friend_models.FriendVerifyModel{
		SendUserID: req.UserID,
		RevUserID:  req.FriendID,
		Message:    req.Verify,
		Source:     req.Source, // 添加来源字段
	}

	err = l.svcCtx.DB.Create(&verifyModel).Error
	if err != nil {
		l.Logger.Errorf("创建好友验证请求失败: %v", err)
		return nil, errors.New("添加好友请求失败")
	}

	l.Logger.Infof("好友请求发送成功: userID=%s, friendID=%s, source=%s", req.UserID, req.FriendID, req.Source)
	return &types.AddFriendRes{}, nil
}
