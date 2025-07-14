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
	// 参数验证
	if req.FriendID == "" {
		return nil, errors.New("好友ID不能为空")
	}

	// 不能查询自己的信息
	if req.UserID == req.FriendID {
		return nil, errors.New("不能查询自己的信息")
	}

	var friend friend_models.FriendModel

	// 通过RPC获取用户信息
	res, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: req.FriendID,
	})
	if err != nil {
		l.Logger.Errorf("获取用户信息失败: friendID=%s, error=%v", req.FriendID, err)
		return nil, errors.New("用户不存在")
	}

	var friendUser user_models.UserModel
	if err := json.Unmarshal(res.Data, &friendUser); err != nil {
		l.Logger.Errorf("解析用户数据失败: %v", err)
		return nil, errors.New("解析用户数据失败")
	}

	// 生成会话Id
	conversationID, err := conversation.GenerateConversation([]string{req.UserID, req.FriendID})
	if err != nil {
		l.Logger.Errorf("生成会话Id失败: %v", err)
		return nil, fmt.Errorf("生成会话Id失败: %v", err)
	}

	// 获取好友备注
	notice := friend.GetUserNotice(req.UserID)

	// 检查是否为好友
	isFriend := friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID)

	response := &types.FriendInfoRes{
		ConversationID: conversationID,
		UserID:         friendUser.UUID,
		Nickname:       friendUser.NickName,
		Avatar:         friendUser.Avatar,
		Abstract:       friendUser.Abstract,
		Notice:         notice,
		IsFriend:       isFriend,
		Email:          friendUser.Email,
	}

	l.Logger.Infof("获取好友信息成功: userID=%s, friendID=%s, isFriend=%v", req.UserID, req.FriendID, isFriend)
	return response, nil
}
