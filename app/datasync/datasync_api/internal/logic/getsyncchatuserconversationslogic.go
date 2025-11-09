package logic

import (
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncChatUserConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的用户会话设置版本
func NewGetSyncChatUserConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncChatUserConversationsLogic {
	return &GetSyncChatUserConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncChatUserConversationsLogic) GetSyncChatUserConversations(req *types.GetSyncChatUserConversationsReq) (resp *types.GetSyncChatUserConversationsRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.Errorf("用户ID为空")
		return nil, errors.New("用户ID不能为空")
	}

	// 获取用户参与的会话列表
	conversationsResp, err := l.svcCtx.ChatRpc.GetUserConversations(l.ctx, &chat_rpc.GetUserConversationsReq{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("获取用户会话列表失败: %v", err)
		return nil, err
	}

	conversationIDs := make([]string, 0, len(conversationsResp.Conversations))
	for _, conv := range conversationsResp.Conversations {
		conversationIDs = append(conversationIDs, conv.ConversationId)
	}

	if len(conversationIDs) == 0 {
		return &types.GetSyncChatUserConversationsRes{
			UserConversationVersions: []types.ChatUserConversationVersionItem{},
			ServerTimestamp:          time.Now().UnixMilli(),
		}, nil
	}

	// 获取变更的用户会话设置版本
	serverTimestamp := time.Now().UnixMilli()

	userConvResp, err := l.svcCtx.ChatRpc.GetUserConversationSettingsListByIds(l.ctx, &chat_rpc.GetUserConversationSettingsListByIdsReq{
		UserId:          userId,
		ConversationIds: conversationIDs,
		Since:           req.Since,
	})
	if err != nil {
		l.Errorf("获取变更的用户会话设置版本失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var userConversationVersions []types.ChatUserConversationVersionItem
	for _, userConv := range userConvResp.UserConversationSettings {
		userConversationVersions = append(userConversationVersions, types.ChatUserConversationVersionItem{
			ConversationID: userConv.ConversationId,
			Version:        userConv.Version,
		})
	}

	return &types.GetSyncChatUserConversationsRes{
		UserConversationVersions: userConversationVersions,
		ServerTimestamp:          serverTimestamp,
	}, nil
}
