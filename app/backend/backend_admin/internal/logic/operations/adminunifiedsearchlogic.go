package operations

import (
	"context"
	"errors"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUnifiedSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUnifiedSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUnifiedSearchLogic {
	return &AdminUnifiedSearchLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminUnifiedSearchLogic) AdminUnifiedSearch(req *types.AdminUnifiedSearchReq) (resp *types.AdminUnifiedSearchRes, err error) {
	keyword := strings.TrimSpace(req.Keyword)
	if keyword == "" {
		return nil, errors.New("请输入检索关键词")
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 30 {
		limit = 30
	}

	resp = &types.AdminUnifiedSearchRes{
		Users:         []types.SearchUserHit{},
		Groups:        []types.SearchGroupHit{},
		Messages:      []types.SearchMessageHit{},
		Conversations: []types.SearchConversationHit{},
	}

	userReq := &user_rpc.ListUsersReq{Keyword: keyword, UserId: keyword, Page: 1, PageSize: int32(limit)}
	if strings.Contains(keyword, "@") {
		userReq.Email = keyword
	}
	userRes, err := l.svcCtx.UserRpc.ListUsers(l.ctx, userReq)
	if err != nil {
		l.Errorf("检索用户失败: %v", err)
	} else {
		for _, u := range userRes.List {
			resp.Users = append(resp.Users, types.SearchUserHit{
				UserID: u.UserId, NickName: u.NickName, Email: u.Email, Status: int(u.Status),
			})
		}
	}

	groupRes, err := l.svcCtx.GroupRpc.ListGroups(l.ctx, &group_rpc.ListGroupsReq{
		Keywords: keyword, GroupId: keyword, Page: 1, PageSize: int32(limit),
	})
	if err != nil {
		l.Errorf("检索群组失败: %v", err)
	} else {
		for _, g := range groupRes.List {
			resp.Groups = append(resp.Groups, types.SearchGroupHit{
				GroupID: g.GroupId, Title: g.Title, Status: int(g.Status),
			})
		}
	}

	var msgList []*chat_rpc.ChatMessageItem
	msgRes, err := l.svcCtx.ChatRpc.ListChatMessages(l.ctx, &chat_rpc.ListChatMessagesReq{
		MessageId: keyword, Page: 1, PageSize: int32(limit), WithContent: false,
	})
	if err != nil {
		l.Errorf("检索消息失败: %v", err)
	} else {
		msgList = msgRes.List
		for _, m := range msgList {
			resp.Messages = append(resp.Messages, types.SearchMessageHit{
				MessageID: m.MessageId, ConversationID: m.ConversationId,
				SendUserID: m.SendUserId, MsgPreview: m.MsgPreview, CreateTime: m.CreatedAt,
			})
		}
	}

	if len(msgList) == 0 {
		convRes, convErr := l.svcCtx.ChatRpc.ListChatMessages(l.ctx, &chat_rpc.ListChatMessagesReq{
			ConversationId: keyword, Page: 1, PageSize: 1,
		})
		if convErr == nil && len(convRes.List) > 0 {
			m := convRes.List[0]
			resp.Conversations = append(resp.Conversations, types.SearchConversationHit{
				ConversationID: m.ConversationId, Title: m.ConversationId, LastMessage: m.MsgPreview,
			})
		}
	}

	return resp, nil
}
