package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/models/ctype"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditMessageLogic {
	return &EditMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditMessageLogic) EditMessage(req *types.EditMessageReq) (resp *types.EditMessageRes, err error) {
	if req.Content == "" {
		return nil, errors.New("消息内容不能为空")
	}

	var msg chat_models.ChatMessage
	if err = l.svcCtx.DB.Where("message_id = ?", req.MessageID).First(&msg).Error; err != nil {
		return nil, errors.New("消息不存在")
	}

	if msg.SendUserID == nil || *msg.SendUserID != req.UserID {
		return nil, errors.New("无权编辑他人消息")
	}
	if msg.Status == 2 {
		return nil, errors.New("消息已撤回，无法编辑")
	}
	if msg.MsgType != ctype.TextMsgType && msg.MsgType != ctype.MarkdownMsgType {
		return nil, errors.New("仅支持编辑文本或 Markdown 消息")
	}
	if time.Since(time.Time(msg.CreatedAt)) > 24*time.Hour {
		return nil, errors.New("超过24小时，无法编辑")
	}

	if msg.Msg == nil {
		msg.Msg = &ctype.Msg{Type: msg.MsgType}
	}
	switch msg.MsgType {
	case ctype.TextMsgType:
		if msg.Msg.TextMsg == nil {
			msg.Msg.TextMsg = &ctype.TextMsg{}
		}
		msg.Msg.TextMsg.Content = req.Content
	case ctype.MarkdownMsgType:
		if msg.Msg.MarkdownMsg == nil {
			msg.Msg.MarkdownMsg = &ctype.MarkdownMsg{}
		}
		msg.Msg.MarkdownMsg.Content = req.Content
	}

	newPreview := msg.MsgPreviewMethod()
	editTime := time.Now()

	if err = l.svcCtx.DB.Model(&msg).Updates(map[string]interface{}{
		"msg":         msg.Msg,
		"msg_preview": newPreview,
		"status":      3,
		"updated_at":  editTime,
	}).Error; err != nil {
		l.Errorf("更新消息失败: messageId=%s, error=%v", req.MessageID, err)
		return nil, errors.New("编辑失败")
	}

	var conversationVersion int64
	var meta chat_models.ChatConversationMeta
	if err := l.svcCtx.DB.Where("conversation_id = ?", msg.ConversationID).First(&meta).Error; err == nil && meta.MaxSeq == msg.Seq {
		conversationVersion = l.svcCtx.VersionGen.GetNextVersion("chat_conversation_metas", "conversation_id", msg.ConversationID)
		if err := l.svcCtx.DB.Model(&meta).Updates(map[string]interface{}{
			"last_message": newPreview,
			"version":      conversationVersion,
			"updated_at":   editTime,
		}).Error; err != nil {
			l.Errorf("更新会话预览失败: conversationId=%s, error=%v", msg.ConversationID, err)
		}
	}

	go l.notifyMessageEdited(req.UserID, msg.ConversationID, msg.Seq, conversationVersion)

	return &types.EditMessageRes{
		Id:        msg.Id,
		MessageID: req.MessageID,
		Content:   req.Content,
		EditTime:  editTime.Format("2006-01-02 15:04:05"),
	}, nil
}

func (l *EditMessageLogic) notifyMessageEdited(senderID, conversationID string, seq, conversationVersion int64) {
	defer func() {
		if r := recover(); r != nil {
			l.Errorf("推送消息编辑通知时发生panic: %v", r)
		}
	}()

	var userConversations []chat_models.ChatUserConversation
	if err := l.svcCtx.DB.Where("conversation_id = ?", conversationID).Find(&userConversations).Error; err != nil {
		l.Errorf("查询会话成员失败: conversationId=%s, error=%v", conversationID, err)
		return
	}

	messagesUpdate := map[string]interface{}{
		"table":          "messages",
		"conversationId": conversationID,
		"data": []map[string]interface{}{
			{"seq": seq},
		},
	}

	tableUpdates := []map[string]interface{}{messagesUpdate}
	if conversationVersion > 0 {
		tableUpdates = append(tableUpdates, map[string]interface{}{
			"table":          "conversations",
			"conversationId": conversationID,
			"data": []map[string]interface{}{
				{"version": conversationVersion},
			},
		})
	}

	for _, uc := range userConversations {
		payload := map[string]interface{}{
			"command":  wsCommandConst.CHAT_MESSAGE,
			"type":     wsTypeConst.ChatConversationMessageReceive,
			"senderId": senderID,
			"targetId": uc.UserID,
			"body": map[string]interface{}{
				"tableUpdates": tableUpdates,
			},
			"conversationId": conversationID,
		}
		if err := l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload); err != nil {
			l.Errorf("推送消息编辑通知失败: target=%s, error=%v", uc.UserID, err)
		}
	}
}
