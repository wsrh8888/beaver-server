package logic

import (
	"context"
	"encoding/json"
	"time"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 聊天数据同步
func NewChatSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatSyncLogic {
	return &ChatSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatSyncLogic) ChatSync(req *types.ChatSyncReq) (resp *types.ChatSyncRes, err error) {
	var chats []chat_models.ChatMessage

	// 过滤掉当前用户主动删除的消息
	deleteSubQuery := l.svcCtx.DB.Model(&chat_models.ChatUserDelete{}).
		Select("message_id").
		Where("user_id = ?", req.UserID)

	// 构建基础查询条件 - 查询指定seq范围的消息（包含起始seq）
	var query = l.svcCtx.DB.Where("seq >= ? AND seq <= ? AND message_id NOT IN (?)", req.FromSeq, req.ToSeq, deleteSubQuery)

	// 如果指定了会话ID，则只同步该会话的消息
	if req.ConversationID != "" {
		query = query.Where("conversation_id = ?", req.ConversationID)
	} else {
		// 如果没有指定会话ID，需要过滤出用户相关的会话
		// 通过子查询获取用户参与的所有会话ID（包括已删除的）
		subQuery := l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
			Select("conversation_id").
			Where("user_id = ?", req.UserID)

		query = query.Where("conversation_id IN (?)", subQuery)
	}

	// 设置默认限制
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	// 查询数据
	err = query.Order("seq ASC").Limit(limit + 1).Find(&chats).Error
	if err != nil {
		l.Errorf("查询聊天数据失败: %v", err)
		return nil, err
	}

	// 判断是否还有更多数据
	hasMore := len(chats) > limit
	if hasMore {
		chats = chats[:limit]
	}

	// 转换为响应格式
	var messages = make([]types.ChatSyncMessage, 0)
	var nextSeq int64 = req.FromSeq

	for _, chat := range chats {
		// 序列化消息内容
		msgJson := ""
		if chat.Msg != nil {
			// 将Msg结构体序列化为JSON字符串
			msgBytes, err := json.Marshal(chat.Msg)
			if err != nil {
				l.Errorf("序列化消息内容失败: %v", err)
				msgJson = chat.MsgPreview // 出错时使用预览
			} else {
				msgJson = string(msgBytes)
			}
		}

		// 原始消息不修改（只增不改原则），撤回状态由同步流中的 WithdrawMsg 指令消息表达
		// isDeleted 保留字段兼容，始终为 false（原始消息不变）
		isDeleted := false

		// 处理发送者ID
		sendUserID := ""
		if chat.SendUserID != nil {
			sendUserID = *chat.SendUserID
		}

		// 对于通知消息，前端可以根据SendUserID是否为空来判断

		messages = append(messages, types.ChatSyncMessage{
			MessageID:        chat.MessageID,
			ConversationID:   chat.ConversationID,
			ConversationType: chat.ConversationType,
			SendUserID:       sendUserID,
			MsgType:          uint32(chat.MsgType),
			MsgPreview:       chat.MsgPreview,
			Msg:              msgJson,
			IsDeleted:        isDeleted,
			Seq:              chat.Seq,
			CreatedAt:        time.Time(chat.CreatedAt).Unix(),
		})

		nextSeq = chat.Seq
	}

	// 如果没有更多数据，nextSeq应该是toSeq+1
	if !hasMore {
		nextSeq = req.ToSeq + 1
	}

	resp = &types.ChatSyncRes{
		Messages: messages,
		HasMore:  hasMore,
		NextSeq:  nextSeq,
	}

	l.Infof("聊天数据同步完成，用户ID: %s, 返回消息数: %d, 还有更多: %v", req.UserID, len(messages), hasMore)
	return resp, nil
}
