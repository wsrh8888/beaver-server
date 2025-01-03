package logic

import (
	"context"
	"fmt"
	"strings"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
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

	pattern := "%" + req.UserId + "_%"  // userId_ 在前
	pattern2 := "%_" + req.UserId + "%" // userId_ 在后
	fmt.Println(req.UserId, "sssssssssssssss")
	err = l.svcCtx.DB.Raw(sql, pattern, pattern2).Scan(&userConversations).Error
	if err != nil {
		return nil, err
	}

	var privateChatIds []string
	var groupChatIds []string

	// 获取会话Ids并根据命名规则区分私聊和群聊
	for _, convo := range userConversations {
		if strings.Contains(convo.ConversationId, "_") {
			privateChatIds = append(privateChatIds, convo.ConversationId)
		} else {
			groupChatIds = append(groupChatIds, convo.ConversationId)
		}
	}

	// 从私聊Id中提取用户Id进行去重
	userIdMap := make(map[string]struct{})
	for _, chatId := range privateChatIds {
		ids := strings.Split(chatId, "_")
		for _, id := range ids {
			if id != req.UserId {
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
		err = l.svcCtx.DB.Where("user_id IN (?)", userIds).Find(&users).Error
		if err != nil {
			return nil, err
		}
	}

	// 构建最近会话列表数据
	userMap := make(map[string]user_models.UserModel)
	for _, user := range users {
		userMap[user.UserId] = user
	}

	var respList []types.RecentChat
	for _, convo := range userConversations {
		var chatInfo types.RecentChat

		chatInfo.MsgPreview = convo.LastMessage
		chatInfo.IsTop = convo.IsPinned
		chatInfo.ConversationId = convo.ConversationId
		chatInfo.CreateAt = convo.CreatedAt.String()
		ids := strings.Split(convo.ConversationId, "_")
		opponentId := ids[0]
		if ids[0] == req.UserId {
			opponentId = ids[1]
		}
		user := userMap[opponentId]
		chatInfo.Nickname = user.NickName
		chatInfo.Avatar = user.Avatar

		respList = append(respList, chatInfo)
	}

	return &types.RecentChatListRes{
		List:  respList,
		Count: int64(len(respList)),
	}, nil
}
