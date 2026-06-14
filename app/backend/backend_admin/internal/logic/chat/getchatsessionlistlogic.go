package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatSessionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// GetChatSessionList 管理后台：按用户列出会话（运营审计入口）。
// admin 职责：校验 userId、UserRpc 组装参与方昵称与会话展示标题。
// RPC 职责：ListConversations 领域查询，不与本 HTTP 接口 1:1。
func NewGetChatSessionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatSessionListLogic {
	return &GetChatSessionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatSessionListLogic) GetChatSessionList(req *types.GetChatSessionListReq) (resp *types.GetChatSessionListRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("请先选择或搜索用户")
	}

	rpcRes, err := l.svcCtx.ChatRpc.ListConversations(l.ctx, &chat_rpc.ListConversationsReq{
		UserId:           req.UserID,
		ConversationType: int32(req.ConversationType),
		Page:             int32(req.Page),
		PageSize:         int32(req.PageSize),
	})
	if err != nil {
		l.Errorf("获取用户会话列表失败: %v", err)
		return nil, err
	}

	users := map[string]*user_rpc.UserInfo{}
	seen := map[string]struct{}{}
	for _, item := range rpcRes.List {
		for _, uid := range item.ParticipantIds {
			if uid == "" {
				continue
			}
			if _, ok := seen[uid]; ok {
				continue
			}
			seen[uid] = struct{}{}
		}
	}
	userIDs := make([]string, 0, len(seen))
	for uid := range seen {
		userIDs = append(userIDs, uid)
	}
	if len(userIDs) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); err == nil && res != nil {
			users = res.UserInfo
		}
	}

	list := make([]types.GetChatSessionListItem, 0, len(rpcRes.List))
	for _, item := range rpcRes.List {
		names := make([]string, 0, len(item.ParticipantIds))
		for _, uid := range item.ParticipantIds {
			name := uid
			if u, ok := users[uid]; ok && u != nil && u.NickName != "" {
				name = u.NickName
			}
			names = append(names, name)
		}

		peerID := ""
		peerName := ""
		title := item.ConversationId
		if item.Type == 1 {
			for _, uid := range item.ParticipantIds {
				if uid != req.UserID {
					peerID = uid
					if u, ok := users[uid]; ok && u != nil {
						peerName = u.NickName
					}
					if peerName == "" {
						peerName = uid
					}
					title = peerName
					break
				}
			}
		} else if item.Type == 2 {
			title = fmt.Sprintf("群聊 · %s", shortenID(item.ConversationId))
		}

		list = append(list, types.GetChatSessionListItem{
			ConversationID:     item.ConversationId,
			ConversationType:   int(item.Type),
			Title:              title,
			PeerUserID:         peerID,
			PeerUserName:       peerName,
			ParticipantIDs:     item.ParticipantIds,
			ParticipantNames:   names,
			LastMessage:        item.LastMessage,
			LastMessageTime:    item.LastMessageTime,
			MessageCount:       item.MessageCount,
		})
	}

	return &types.GetChatSessionListRes{List: list, Total: rpcRes.Total}, nil
}

func shortenID(id string) string {
	if len(id) <= 12 {
		return id
	}
	return id[:6] + "…" + id[len(id)-4:]
}
