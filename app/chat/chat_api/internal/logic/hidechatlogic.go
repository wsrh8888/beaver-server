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

type HideChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 隐藏/显示会话
func NewHideChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HideChatLogic {
	return &HideChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HideChatLogic) HideChat(req *types.HideChatReq) (resp *types.HideChatRes, err error) {
	resp = &types.HideChatRes{}

	// 获取下一个版本号
	version := l.svcCtx.VersionGen.GetNextVersion("chat_user_conversations", "user_id", req.UserID)

	// 更新会话隐藏状态和版本号
	err = l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
		Where("user_id = ? AND conversation_id = ?", req.UserID, req.ConversationID).
		Updates(map[string]interface{}{
			"is_hidden": req.IsHidden,
			"version":   version,
		}).Error
	if err != nil {
		l.Logger.Errorf("hide chat update failed: %v", err)
		return nil, err
	}

	// 发送WS通知给自己（更新本地数据）
	go func() {
		l.notifyHiddenUpdate(req.ConversationID, req.UserID, version)
	}()

	return resp, nil
}

// 发送隐藏状态更新通知
func (l *HideChatLogic) notifyHiddenUpdate(conversationId, userId string, version int64) {
	defer func() {
		if r := recover(); r != nil {
			l.Logger.Errorf("发送隐藏通知时发生panic: %v", r)
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

	l.Logger.Infof("发送隐藏状态更新通知: user=%s, conversation=%s, version=%d",
		userId, conversationId, version)
}
