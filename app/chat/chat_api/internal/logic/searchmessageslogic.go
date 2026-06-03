package logic

import (
	"context"
	"errors"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/list_query"
	"beaver/common/models"
	"beaver/common/models/ctype"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchMessagesLogic {
	return &SearchMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchMessagesLogic) SearchMessages(req *types.SearchMessagesReq) (*types.SearchMessagesRes, error) {
	if req.Keyword == "" {
		return nil, errors.New("keyword不能为空")
	}

	keyword := "%" + req.Keyword + "%"

	deleteSubQuery := l.svcCtx.DB.Model(&chat_models.ChatUserDelete{}).
		Select("message_id").
		Where("user_id = ?", req.UserID)

	conversationSubQuery := l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
		Select("conversation_id").
		Where("user_id = ?", req.UserID)

	query := l.svcCtx.DB.Where(
		"msg_preview LIKE ? AND message_id NOT IN (?) AND conversation_id IN (?)",
		keyword, deleteSubQuery, conversationSubQuery,
	).Where("msg_type IN ?", []ctype.MsgType{ctype.TextMsgType, ctype.MarkdownMsgType})

	if req.ConversationID != "" {
		var count int64
		l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
			Where("user_id = ? AND conversation_id = ?", req.UserID, req.ConversationID).
			Count(&count)
		if count == 0 {
			return nil, errors.New("无权搜索该会话")
		}
		query = query.Where("conversation_id = ?", req.ConversationID)
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	chatMessages, total, err := list_query.ListQuery(l.svcCtx.DB, chat_models.ChatMessage{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: limit,
			Sort:  "created_at desc",
		},
		Where: query,
	})
	if err != nil {
		return nil, err
	}

	userIds := make([]string, 0)
	userIdSet := make(map[string]bool)
	for _, chat := range chatMessages {
		if chat.SendUserID != nil && *chat.SendUserID != "" && !userIdSet[*chat.SendUserID] {
			userIds = append(userIds, *chat.SendUserID)
			userIdSet[*chat.SendUserID] = true
		}
	}

	userInfoMap := make(map[string]types.Sender)
	if len(userIds) > 0 {
		userListResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
			UserIdList: userIds,
		})
		if err != nil {
			l.Errorf("批量获取用户信息失败: %v", err)
		} else {
			for userId, userInfo := range userListResp.UserInfo {
				userInfoMap[userId] = types.Sender{
					UserID:   userId,
					NickName: userInfo.NickName,
					Avatar:   userInfo.Avatar,
				}
			}
		}
	}

	list := make([]types.Message, 0, len(chatMessages))
	for _, chat := range chatMessages {
		var msg types.Msg
		if chat.Msg != nil {
			if err := convertCtypeMsgToTypesMsg(*chat.Msg, &msg); err != nil {
				return nil, err
			}
		}

		sendUserID := ""
		if chat.SendUserID != nil {
			sendUserID = *chat.SendUserID
		}

		sender := types.Sender{UserID: sendUserID, NickName: "未知用户"}
		if sendUserID == "" {
			sender = types.Sender{NickName: "通知消息"}
		} else if info, ok := userInfoMap[sendUserID]; ok {
			sender = info
		}

		list = append(list, types.Message{
			Id:               chat.Id,
			MessageID:        chat.MessageID,
			ConversationID:   chat.ConversationID,
			ConversationType: chat.ConversationType,
			Sender:           sender,
			CreatedAt:        chat.CreatedAt.String(),
			Msg:              msg,
			Seq:              chat.Seq,
		})
	}

	return &types.SearchMessagesRes{
		Count: total,
		List:  list,
	}, nil
}
