package logic

import (
	"context"
	"fmt"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/common/list_query"
	"beaver/common/models"
	"beaver/utils/conversation"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListRes, err error) {
	friends, _, _ := list_query.ListQuery(l.svcCtx.DB, friend_models.FriendModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Where:   l.svcCtx.DB.Where("(send_user_id = ? OR rev_user_id = ?) AND is_deleted = ?", req.UserID, req.UserID, false),
		Preload: []string{"SendUserModel", "RevUserModel"},
	})

	var list []types.FriendInfoRes
	for _, friendUser := range friends {
		info := types.FriendInfoRes{}

		if friendUser.SendUserID == req.UserID {
			conversationID, err := conversation.GenerateConversation([]string{req.UserID, friendUser.RevUserModel.UUID})
			if err != nil {
				return nil, fmt.Errorf("生成会话Id失败: %v", err)
			}
			// 我是发起方
			info = types.FriendInfoRes{
				UserID:         friendUser.RevUserModel.UUID,
				Nickname:       friendUser.RevUserModel.NickName,
				Avatar:         friendUser.RevUserModel.Avatar,
				Abstract:       friendUser.RevUserModel.Abstract,
				Notice:         friendUser.SendUserNotice,
				ConversationID: conversationID,
				Phone:          friendUser.RevUserModel.Phone,
			}
		}
		if friendUser.RevUserID == req.UserID {
			conversationID, err := conversation.GenerateConversation([]string{req.UserID, friendUser.SendUserModel.UUID})
			if err != nil {
				return nil, fmt.Errorf("生成会话Id失败: %v", err)
			}
			// 我是接收方
			info = types.FriendInfoRes{
				UserID:         friendUser.SendUserModel.UUID,
				Nickname:       friendUser.SendUserModel.NickName,
				Avatar:         friendUser.SendUserModel.Avatar,
				Abstract:       friendUser.SendUserModel.Abstract,
				Notice:         friendUser.RevUserNotice,
				ConversationID: conversationID,
				Phone:          friendUser.SendUserModel.Phone,
			}
		}
		list = append(list, info)
	}

	return &types.FriendListRes{
		List: list,
	}, nil
}
