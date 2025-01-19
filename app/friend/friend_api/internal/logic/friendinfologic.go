package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/conversation"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendInfoLogic {
	return &FriendInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendInfoLogic) FriendInfo(req *types.FriendInfoReq) (resp *types.FriendInfoRes, err error) {

	var friend friend_models.FriendModel

	res, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: req.FriendID,
	})
	if err != nil {
		return nil, errors.New(err.Error())
	}
	var user user_models.UserModel
	json.Unmarshal([]byte(res.Data), &user)

	var friendUser user_models.UserModel
	json.Unmarshal(res.Data, &friendUser)
	conversationID, err := conversation.GenerateConversation([]string{req.UserID, req.FriendID})
	if err != nil {
		return nil, fmt.Errorf("生成会话Id失败: %v", err)
	}
	response := &types.FriendInfoRes{
		ConversationID: conversationID,
		UserID:         friendUser.UUID,
		Nickname:       friendUser.NickName,
		Avatar:         friendUser.Avatar,
		Abstract:       friendUser.Abstract,
		Notice:         friend.GetUserNotice(req.UserID),
		IsFriend:       friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID),
		Phone:          friendUser.Phone,
	}

	return response, nil
}
