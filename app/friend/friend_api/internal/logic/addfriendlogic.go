package logic

import (
	"context"
	"errors"
	"fmt"

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
	if req.UserId == req.FriendId {
		return nil, errors.New("不能添加自己为好友")
	}
	var friend friend_models.FriendModel
	if friend.IsFriend(l.svcCtx.DB, req.UserId, req.FriendId) {
		return nil, errors.New("已经是好友了")
	}
	var userInfo user_models.UserModel

	err = l.svcCtx.DB.Take(&userInfo, "user_id = ?", req.FriendId).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	err = l.svcCtx.DB.Take(&friend_models.FriendVerifyModel{}, "(send_user_id = ? AND rev_user_id = ? AND rev_status = 0) OR (send_user_id = ? AND rev_user_id = ? AND rev_status = 0)", req.UserId, req.FriendId, req.FriendId, req.UserId).Error
	if err == nil {
		fmt.Println("当前已经有好友请求")
		return nil, nil
	}

	resp = new(types.AddFriendRes)
	var verifyModel = friend_models.FriendVerifyModel{
		SendUserId: req.UserId,
		RevUserId:  req.FriendId,
		Message:    req.Verify,
	}

	err = l.svcCtx.DB.Create(&verifyModel).Error
	if err != nil {
		fmt.Println("添加失败", err.Error())
		return nil, errors.New("添加失败")
	}
	return
}
