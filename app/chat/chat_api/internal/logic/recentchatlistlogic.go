package logic

import (
	"context"
	"fmt"
	"strings"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/group/group_models"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecentChatListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecentChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecentChatListLogic {
	return &RecentChatListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecentChatListLogic) RecentChatList(req *types.RecentChatListReq) (resp *types.RecentChatListRes, err error) {
	var userConversations []chat_models.ChatUserConversationModel

	// 获取该用户的所有未删除会话及其最新一条消息，并按更新时间排序
	sql := `
	SELECT c.*
	FROM chat_user_conversation_models c
	INNER JOIN (
					SELECT conversation_id, MAX(updated_at) AS max_updated_at
					FROM chat_user_conversation_models
					WHERE is_deleted = false 
						AND (conversation_id LIKE ? OR conversation_id LIKE ?)
					GROUP BY conversation_id
	) AS latest
	ON c.conversation_id = latest.conversation_id AND c.updated_at = latest.max_updated_at
	ORDER BY latest.max_updated_at DESC`

	pattern := "%" + req.UserID + "_%"  // userId_ 在前
	pattern2 := "%_" + req.UserID + "%" // userId_ 在后
	fmt.Println(req.UserID, "sssssssssssssss")
	err = l.svcCtx.DB.Raw(sql, pattern, pattern2).Scan(&userConversations).Error
	if err != nil {
		return nil, err
	}

	var privateChatIds []string
	var groupChatIds []string

	// 获取会话Ids并根据命名规则区分私聊和群聊
	for _, convo := range userConversations {
		if strings.Contains(convo.ConversationID, "_") {
			privateChatIds = append(privateChatIds, convo.ConversationID)
		} else {
			groupChatIds = append(groupChatIds, convo.ConversationID)
		}
	}

	// 从私聊Id中提取用户Id进行去重
	userIdMap := make(map[string]struct{})
	for _, chatID := range privateChatIds {
		ids := strings.Split(chatID, "_")
		for _, id := range ids {
			if id != req.UserID {
				userIdMap[id] = struct{}{}
			}
		}
	}

	// 将map转为slice
	var userIds []string
	for id := range userIdMap {
		userIds = append(userIds, id)
	}

	// 批量查询用户信息
	var users []user_models.UserModel
	if len(userIds) > 0 {
		err = l.svcCtx.DB.Where("uuid IN (?)", userIds).Find(&users).Error
		if err != nil {
			return nil, err
		}
	}

	// 批量查询群组信息
	var groups []group_models.GroupModel
	if len(groupChatIds) > 0 {
		err = l.svcCtx.DB.Where("uuid IN (?)", groupChatIds).Find(&groups).Error
		if err != nil {
			return nil, err
		}
	}

	// 构建最近会话列表数据
	userMap := make(map[string]user_models.UserModel)
	for _, user := range users {
		userMap[user.UUID] = user
	}

	groupMap := make(map[string]group_models.GroupModel)
	for _, group := range groups {
		groupMap[group.UUID] = group
	}

	var respList []types.RecentChat
	for _, convo := range userConversations {
		var chatInfo types.RecentChat

		chatInfo.MsgPreview = convo.LastMessage
		chatInfo.IsTop = convo.IsPinned
		chatInfo.ConversationID = convo.ConversationID
		chatInfo.CreateAt = convo.CreatedAt.String()
		if strings.Contains(convo.ConversationID, "_") { // 私聊
			ids := strings.Split(convo.ConversationID, "_")
			opponentID := ids[0]
			if ids[0] == req.UserID {
				opponentID = ids[1]
			}
			user := userMap[opponentID]
			chatInfo.Nickname = user.NickName
			chatInfo.Avatar = user.Avatar
		} else { // 群聊
			group := groupMap[convo.ConversationID]
			chatInfo.Nickname = group.Title
		}

		respList = append(respList, chatInfo)
	}

	return &types.RecentChatListRes{
		List:  respList,
		Count: int64(len(respList)),
	}, nil
}
