package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_rpc/types/user_rpc"
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
	// 参数验证
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// 查询好友列表
	friends, _, _ := list_query.ListQuery(l.svcCtx.DB, friend_models.FriendModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Where: l.svcCtx.DB.Where("(send_user_id = ? OR rev_user_id = ?) AND is_deleted = ?", req.UserID, req.UserID, false),
		// 移除Preload，微服务架构中通过RPC获取用户信息
	})

	// 收集需要获取用户信息的UserID列表
	var userIds []string
	userIdSet := make(map[string]bool)
	for _, friendUser := range friends {
		if friendUser.SendUserID != "" && !userIdSet[friendUser.SendUserID] {
			userIds = append(userIds, friendUser.SendUserID)
			userIdSet[friendUser.SendUserID] = true
		}
		if friendUser.RevUserID != "" && !userIdSet[friendUser.RevUserID] {
			userIds = append(userIds, friendUser.RevUserID)
			userIdSet[friendUser.RevUserID] = true
		}
	}

	// 批量获取用户信息
	userInfoMap := make(map[string]*user_rpc.UserInfo)
	if len(userIds) > 0 {
		userListResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
			UserIdList: userIds,
		})
		if err != nil {
			l.Logger.Errorf("批量获取用户信息失败: %v", err)
			// 不返回错误，继续处理，为没有用户信息的设置默认值
		} else {
			userInfoMap = userListResp.UserInfo
		}
	}

	var list []types.FriendInfoRes
	for _, friendUser := range friends {
		var targetUserID string
		var notice string

		// 确定目标用户ID和备注信息
		if friendUser.SendUserID == req.UserID {
			// 我是发起方，目标用户是接收方
			targetUserID = friendUser.RevUserID
			notice = friendUser.SendUserNotice
		} else if friendUser.RevUserID == req.UserID {
			// 我是接收方，目标用户是发起方
			targetUserID = friendUser.SendUserID
			notice = friendUser.RevUserNotice
		} else {
			// 这种情况理论上不应该发生，跳过
			continue
		}

		// 获取用户信息
		var nickname, avatar, abstract, email string
		if userInfo, exists := userInfoMap[targetUserID]; exists && userInfo != nil {
			nickname = userInfo.NickName
			avatar = userInfo.Avatar
			email = userInfo.Email

			// 获取完整的用户信息（包括Abstract字段）
			userDetailResp, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
				UserID: targetUserID,
			})
			if err == nil {
				abstract = userDetailResp.UserInfo.Abstract
			} else {
				abstract = ""
			}
		} else {
			nickname = "未知用户"
			avatar = ""
			abstract = ""
			email = ""
		}

		// 生成会话Id
		conversationID, err := conversation.GenerateConversation([]string{req.UserID, targetUserID})
		if err != nil {
			l.Logger.Errorf("生成会话Id失败: userID=%s, targetID=%s, error=%v", req.UserID, targetUserID, err)
			return nil, fmt.Errorf("生成会话Id失败: %v", err)
		}

		// 构造好友信息
		info := types.FriendInfoRes{
			UserID:         targetUserID,
			Nickname:       nickname,
			Avatar:         avatar,
			Abstract:       abstract,
			Notice:         notice,
			ConversationID: conversationID,
			Email:          email,
		}

		list = append(list, info)
	}

	l.Logger.Infof("获取好友列表成功: userID=%s, count=%d", req.UserID, len(list))
	return &types.FriendListRes{
		List: list,
	}, nil
}
