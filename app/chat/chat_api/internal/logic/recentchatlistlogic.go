package logic

import (
	"context"
	"strings"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/group/group_models"

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
		if strings.Contains(convo.ConversationID, "_") {
			privateChatIds = append(privateChatIds, convo.ConversationID)
		} else {
			groupChatIds = append(groupChatIds, convo.ConversationID)
		}
	}

	// 从私聊Id中提取用户Id进行去重
	// 会话ID格式: private_A_B
	userIdMap := make(map[string]struct{})
	for _, chatID := range privateChatIds {
		ids := strings.Split(chatID, "_")
		// 跳过第一个元素 "private"，取后面两个用户ID
		if len(ids) >= 3 {
			for i := 1; i < len(ids); i++ {
				id := ids[i]
				if id != req.UserID && id != "private" {
					userIdMap[id] = struct{}{}
				}
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

	// 批量查询群组信息
	var groups []group_models.GroupModel
	if len(groupChatIds) > 0 {
		err = l.svcCtx.DB.Where("group_id IN (?)", groupChatIds).Find(&groups).Error
		if err != nil {
			return nil, err
		}
	}

	// 构建数据映射
	groupMap := make(map[string]group_models.GroupModel)
	for _, group := range groups {
		groupMap[group.GroupID] = group
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
		chatInfo.UpdateAt = convo.UpdatedAt.String()
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
			group := groupMap[convo.ConversationID]
			chatInfo.NickName = group.Title
			chatInfo.Avatar = group.Avatar
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
