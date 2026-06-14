package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

const chatMessageStatusDeleted int32 = 4

type GetChatMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChatMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatMessageListLogic {
	return &GetChatMessageListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetChatMessageListLogic) GetChatMessageList(req *types.GetChatMessageListReq) (resp *types.GetChatMessageListRes, err error) {
	rpcReq := &chat_rpc.ListChatMessagesReq{
		Page:           int32(req.Page),
		PageSize:       int32(req.PageSize),
		ConversationId: req.ConversationID,
		SendUserId:     req.SendUserID,
		MsgType:        int32(req.MsgType),
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		WithContent:    req.WithContent,
	}
	if req.Order == 1 {
		rpcReq.Order = 1
	} else if req.Order == 2 {
		rpcReq.Order = 2
	}
	if req.IsDeleted {
		rpcReq.Status = chatMessageStatusDeleted
	}

	rpcRes, err := l.svcCtx.ChatRpc.ListChatMessages(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("获取聊天消息列表失败: %v", err)
		return nil, err
	}

	users := map[string]*user_rpc.UserInfo{}
	seen := map[string]struct{}{}
	userIDs := make([]string, 0, len(rpcRes.List))
	for _, m := range rpcRes.List {
		if m.SendUserId == "" {
			continue
		}
		if _, ok := seen[m.SendUserId]; ok {
			continue
		}
		seen[m.SendUserId] = struct{}{}
		userIDs = append(userIDs, m.SendUserId)
	}
	if len(userIDs) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); err == nil && res != nil {
			users = res.UserInfo
		}
	}

	list := make([]types.GetChatMessageListItem, 0, len(rpcRes.List))
	for _, m := range rpcRes.List {
		sendName := ""
		if u, ok := users[m.SendUserId]; ok && u != nil {
			sendName = u.NickName
		}
		list = append(list, types.GetChatMessageListItem{
			Id:               m.MessageId,
			MessageID:        m.MessageId,
			ConversationID:   m.ConversationId,
			ConversationType: int(m.ConversationType),
			SendUserID:       m.SendUserId,
			SendUserName:     sendName,
			MsgType:          int(m.MsgType),
			MsgPreview:       m.MsgPreview,
			MsgContent:       m.MsgContent,
			IsDeleted:        m.Status == chatMessageStatusDeleted,
			CreateTime:       m.CreatedAt,
			UpdateTime:     m.UpdatedAt,
		})
	}
	return &types.GetChatMessageListRes{List: list, Total: rpcRes.Total}, nil
}
