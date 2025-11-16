package logic

import (
	"context"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type PinnedChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPinnedChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PinnedChatLogic {
	return &PinnedChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PinnedChatLogic) PinnedChat(req *types.PinnedChatReq) (resp *types.PinnedChatRes, err error) {
	resp = &types.PinnedChatRes{}

	// 获取下一个版本号
	version := l.svcCtx.VersionGen.GetNextVersion("chat_user_conversations", "user_id", req.UserID)

	// 更新会话置顶状态和版本号
	err = l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
		Where("user_id = ? AND conversation_id = ?", req.UserID, req.ConversationID).
		Updates(map[string]interface{}{
			"is_pinned": req.IsPinned,
			"version":   version,
		}).Error
	if err != nil {
		l.Logger.Errorf("pinned chat update failed: %v", err)
		return nil, err
	}

	// 发送WS通知给自己（更新本地数据）
	go func() {
		l.notifyPinnedUpdate(req.ConversationID, req.UserID, version)
	}()

	return resp, nil
}

// 发送置顶状态更新通知
func (l *PinnedChatLogic) notifyPinnedUpdate(conversationId, userId string, version int64) {
	defer func() {
		if r := recover(); r != nil {
			l.Logger.Errorf("发送置顶通知时发生panic: %v", r)
		}
	}()

	// 构建用户会话表更新数据
	userConversationsUpdate := map[string]interface{}{
		"table":          "user_conversations",
		"userId":         userId,
		"conversationId": conversationId,
		"data": []map[string]interface{}{
			{
				"version": int32(version),
			},
		},
	}

	// 发送给自己
	tableUpdates := []map[string]interface{}{userConversationsUpdate}
	messageType := wsTypeConst.ChatUserConversationReceive

	ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.CHAT_MESSAGE, messageType, userId, userId, map[string]interface{}{
		"tableUpdates": tableUpdates,
	}, conversationId)

	l.Logger.Infof("发送置顶状态更新通知: user=%s, conversation=%s, version=%d",
		userId, conversationId, version)
}
