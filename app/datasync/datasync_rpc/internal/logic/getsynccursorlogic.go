package logic

import (
	"context"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/datasync/datasync_rpc/internal/svc"
	"beaver/app/datasync/datasync_rpc/types/types/datasync_rpc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncCursorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSyncCursorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncCursorLogic {
	return &GetSyncCursorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取同步游标
func (l *GetSyncCursorLogic) GetSyncCursor(in *datasync_rpc.GetSyncCursorReq) (*datasync_rpc.GetSyncCursorRes, error) {
	var serverLast int64

	// 根据数据类型获取服务端的最新游标/版本（来自各业务RPC）
	switch in.DataType {
	case "chat_messages":
		// 当前实现采用全局拉取，不按会话维度；固定查询全局最新 seq
		serverLast = l.getLatestChatSeq("")
	case "chat_datasync":
		serverLast = l.getLatestDatasyncVersion()
	case "chat_conversation_settings":
		serverLast = l.getLatestConversationSettingVersion(in.UserId)
	case "friends":
		serverLast = l.getLatestFriendVersion(in.UserId)
	case "friend_verify":
		serverLast = l.getLatestFriendVerifyVersion(in.UserId)
	case "groups":
		serverLast = l.getLatestGroupVersion(in.UserId)
	case "group_members":
		serverLast = l.getLatestGroupMemberVersion(in.UserId)
	case "group_join_requests":
		serverLast = l.getLatestGroupJoinRequestVersion(in.UserId)
	case "users":
		serverLast = l.getLatestUserVersion(in.UserId)
	default:
		serverLast = 0
	}

	return &datasync_rpc.GetSyncCursorRes{
		ServerLatest: serverLast,
	}, nil
}

// getLatestChatSeq 获取最新聊天消息序列号
func (l *GetSyncCursorLogic) getLatestChatSeq(conversationID string) int64 {
	resp, err := l.svcCtx.ChatRpc.GetMessageSeq(l.ctx, &chat_rpc.GetLatestSeqReq{
		ConversationId: conversationID,
	})
	if err != nil {
		l.Errorf("调用chat_rpc获取最新序列号失败: %v", err)
		return 0
	}
	return resp.LatestSeq
}

// getLatestDatasyncVersion 获取最新数据同步表版本号
func (l *GetSyncCursorLogic) getLatestDatasyncVersion() int64 {
	resp, err := l.svcCtx.ChatRpc.GetConversationVersion(l.ctx, &chat_rpc.GetConversationVersionReq{})
	if err != nil {
		l.Errorf("调用chat_rpc获取最新数据同步表版本号失败: %v", err)
		return 0
	}
	return resp.LatestVersion
}

// getLatestConversationSettingVersion 获取最新会话设置版本号
func (l *GetSyncCursorLogic) getLatestConversationSettingVersion(userID string) int64 {
	resp, err := l.svcCtx.ChatRpc.GetConversationSettingVersion(l.ctx, &chat_rpc.GetConversationSettingVersionReq{
		UserId: userID,
	})
	if err != nil {
		l.Errorf("调用chat_rpc获取最新会话设置版本号失败: %v", err)
		return 0
	}
	return resp.LatestVersion
}

// getLatestFriendVersion 获取最新好友版本号
func (l *GetSyncCursorLogic) getLatestFriendVersion(userID string) int64 {
	resp, err := l.svcCtx.FriendRpc.GetFriendVersion(l.ctx, &friend_rpc.GetFriendVersionReq{
		UserId: userID,
	})
	if err != nil {
		l.Errorf("调用friend_rpc获取最新好友版本号失败: %v", err)
		return 0
	}
	return resp.LatestVersion
}

// getLatestFriendVerifyVersion 获取最新好友验证版本号
func (l *GetSyncCursorLogic) getLatestFriendVerifyVersion(userID string) int64 {
	resp, err := l.svcCtx.FriendRpc.GetFriendVerifyVersion(l.ctx, &friend_rpc.GetFriendVerifyVersionReq{
		UserId: userID,
	})
	if err != nil {
		l.Errorf("调用friend_rpc获取最新好友验证版本号失败: %v", err)
		return 0
	}
	return resp.LatestVersion
}

// getLatestGroupVersion 获取最新群组版本号
func (l *GetSyncCursorLogic) getLatestGroupVersion(userID string) int64 {
	resp, err := l.svcCtx.GroupRpc.GetGroupVersion(l.ctx, &group_rpc.GetGroupVersionReq{
		UserId: userID,
	})
	if err != nil {
		l.Errorf("调用group_rpc获取最新群组版本号失败: %v", err)
		return 0
	}
	return resp.LatestVersion
}

// getLatestGroupMemberVersion 获取最新群成员版本号
func (l *GetSyncCursorLogic) getLatestGroupMemberVersion(userID string) int64 {
	resp, err := l.svcCtx.GroupRpc.GetGroupMemberVersion(l.ctx, &group_rpc.GetGroupMemberVersionReq{
		UserId: userID,
	})
	if err != nil {
		l.Errorf("调用group_rpc获取最新群成员版本号失败: %v", err)
		return 0
	}
	return resp.LatestVersion
}

// getLatestGroupJoinRequestVersion 获取最新群组申请版本号
func (l *GetSyncCursorLogic) getLatestGroupJoinRequestVersion(userID string) int64 {
	resp, err := l.svcCtx.GroupRpc.GetGroupJoinRequestVersion(l.ctx, &group_rpc.GetGroupJoinRequestVersionReq{
		UserId: userID,
	})
	if err != nil {
		l.Errorf("调用group_rpc获取最新群组申请版本号失败: %v", err)
		return 0
	}
	return resp.LatestVersion
}

// getLatestUserVersion 获取最新用户版本号
func (l *GetSyncCursorLogic) getLatestUserVersion(userID string) int64 {
	resp, err := l.svcCtx.UserRpc.GetUserVersion(l.ctx, &user_rpc.GetUserVersionReq{
		UserId: userID,
	})
	if err != nil {
		l.Errorf("调用user_rpc获取最新用户版本号失败: %v", err)
		return 0
	}
	return resp.LatestVersion
}
