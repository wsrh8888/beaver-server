package logic

import (
	"context"
	"strings"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/utils/conversation"

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
	var userConversations []chat_models.ChatUserConversation
	var conversationMetas []chat_models.ChatConversationMeta

	// 获取该用户的所有未隐藏会话
	err = l.svcCtx.DB.Where("user_id = ? AND is_hidden = ?", req.UserID, false).
		Order("is_pinned DESC, updated_at DESC").
		Find(&userConversations).Error
	if err != nil {
		return nil, err
	}

	// 获取对应的会话元数据（包含LastMessage）
	conversationIds := make([]string, len(userConversations))
	for i, conv := range userConversations {
		conversationIds[i] = conv.ConversationID
	}

	if len(conversationIds) > 0 {
		err = l.svcCtx.DB.Where("conversation_id IN (?)", conversationIds).
			Find(&conversationMetas).Error
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	var privateChatIds []string
	var groupChatIds []string

	// 获取会话Ids并根据命名规则区分私聊和群聊
	for _, convo := range userConversations {
		conversationType := conversation.GetConversationType(convo.ConversationID)
		if conversationType == 1 { // 私聊
			privateChatIds = append(privateChatIds, convo.ConversationID)
		} else if conversationType == 2 { // 群聊
			groupChatIds = append(groupChatIds, convo.ConversationID)
		}
	}

	// 从私聊Id中提取用户Id进行去重
	userIdMap := make(map[string]struct{})
	for _, chatID := range privateChatIds {
		_, userIds := conversation.ParseConversationWithType(chatID)
		for _, userId := range userIds {
			if userId != req.UserID {
				userIdMap[userId] = struct{}{}
			}
		}
	}

	// 将map转为slice
	var userIds []string
	for id := range userIdMap {
		userIds = append(userIds, id)
	}

	// 通过FriendRpc获取好友详细信息（包含用户基础信息和备注）
	var friendDetails []*friend_rpc.FriendDetailItem
	if len(userIds) > 0 {
		friendDetailRes, err := l.svcCtx.FriendRpc.GetFriendDetail(l.ctx, &friend_rpc.GetFriendDetailReq{
			UserId:    req.UserID,
			FriendIds: userIds,
		})
		if err != nil {
			logx.Errorf("获取好友详情失败: %v", err)
			return nil, err
		}
		friendDetails = friendDetailRes.Friends
	}

	// 通过RPC批量获取群组信息
	var groupDetails []*group_rpc.GroupListById
	if len(groupChatIds) > 0 {
		// 从群聊会话ID中提取实际的群组ID（去掉group_前缀）
		var actualGroupIds []string
		for _, groupChatId := range groupChatIds {
			_, userIds := conversation.ParseConversationWithType(groupChatId)
			if len(userIds) > 0 {
				actualGroupIds = append(actualGroupIds, userIds[0])
			}
		}

		if len(actualGroupIds) > 0 {
			groupRes, err := l.svcCtx.GroupRpc.GetGroupsListByIds(l.ctx, &group_rpc.GetGroupsListByIdsReq{
				GroupIDs: actualGroupIds,
			})
			if err != nil {
				logx.Errorf("获取群组详情失败: %v", err)
				return nil, err
			}
			groupDetails = groupRes.Groups
		}
	}

	// 构建群组数据映射 (key为完整的会话ID)
	groupMap := make(map[string]*group_rpc.GroupListById)
	for i, groupChatId := range groupChatIds {
		if i < len(groupDetails) {
			groupMap[groupChatId] = groupDetails[i]
		}
	}

	// 构建好友详细信息映射
	friendDetailMap := make(map[string]*friend_rpc.FriendDetailItem)
	for _, friendDetail := range friendDetails {
		friendDetailMap[friendDetail.UserId] = friendDetail
	}

	conversationMetaMap := make(map[string]chat_models.ChatConversationMeta)
	for _, meta := range conversationMetas {
		conversationMetaMap[meta.ConversationID] = meta
	}

	var respList []types.ConversationInfoRes
	for _, convo := range userConversations {
		var chatInfo types.ConversationInfoRes

		// 从会话元数据中获取最后消息
		if meta, exists := conversationMetaMap[convo.ConversationID]; exists {
			chatInfo.MsgPreview = meta.LastMessage
		} else {
			chatInfo.MsgPreview = ""
		}

		chatInfo.IsTop = convo.IsPinned
		chatInfo.ConversationID = convo.ConversationID
		chatInfo.UpdatedAt = convo.UpdatedAt.String()
		if strings.HasPrefix(convo.ConversationID, "private_") { // 私聊
			ids := strings.Split(convo.ConversationID, "_")
			// ids格式: ["private", "A", "B"]
			var opponentID string
			if ids[1] == req.UserID {
				opponentID = ids[2]
			} else {
				opponentID = ids[1]
			}

			// 从好友详细信息中获取用户信息和备注
			if friendDetail, exists := friendDetailMap[opponentID]; exists {
				chatInfo.NickName = friendDetail.NickName
				chatInfo.Avatar = friendDetail.Avatar
				chatInfo.Notice = friendDetail.Notice
			} else {
				chatInfo.NickName = "未知用户"
				chatInfo.Avatar = ""
				chatInfo.Notice = ""
			}
			chatInfo.ChatType = 1
		} else { // 群聊
			group, exists := groupMap[convo.ConversationID]
			if exists {
				chatInfo.NickName = group.Name
				chatInfo.Avatar = group.Avatar
			} else {
				chatInfo.NickName = "未知群聊"
				chatInfo.Avatar = ""
			}
			chatInfo.ChatType = 2
			chatInfo.Notice = "" // 群聊暂时没有备注功能
		}

		respList = append(respList, chatInfo)
	}

	return &types.RecentChatListRes{
		List:  respList,
		Count: int64(len(respList)),
	}, nil
}
